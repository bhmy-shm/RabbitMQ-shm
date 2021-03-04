package CheckA

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Lib"
	"RabbitMQ/Middle"
	"RabbitMQ/Model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

/*
A公司的无脑补偿机制。无脑补偿机制主要用来监测日志记录表
如果日志记录表中字段发生改变，根据不同的变化来判断消息转账记录的 延迟过期、失败、重发
*/

var (
	//补偿机制(一)：超时过期，将超时时间>30,且status<>2的 status字段修改成2。
	Timeout = `update translog set STATUS=2 where TIMESTAMPDIFF(SECOND,updatetime,now()) > 30 and STATUS =0`

	//补偿机制(二)：超时后还钱(事务)，一直监测数据库字段，如果有 status=2，且isback=0 的 就直接执行还钱。
	BackSql = "select tid, `from`,money from translog  where status=2 and isback=0 limit 10"
	replay = "update usermoney set user_money=user_money + ? where user_name = ?"
	isback = "update translog set isback=1 where tid = ? "
	Backlock = false

	//补偿机制(三)：转账重发，查询status=0，且超时时间小于10秒的，重发一次
	Retrans = "select * from translog where status=0 and TIMESTAMPDIFF(SECOND,updatetime,now())<=10 "

	//消息队列发送成功，将status=1
	ConfirmSql = "update a.translog set status=1 where status=0"
)

//补偿机制一超时过期的封装函数
func TransTimeOut(){
	_,err := AppInit.GetDB().Exec(Timeout)
	Middle.CheckError(err,"超时过期")
}


//补偿机制二，超时后还钱（事务）
func TransBackMoney(){
	if Backlock {
		log.Println("locked.return.......")
		return
	}

	//1.开启事务
	tx,err := AppInit.GetDB().BeginTxx(context.Background(),nil)
	if err != nil {
		log.Println("开启事务失败：",err)
		return
	}

	Backlock = true

	//2.在事务中执行sql语句，拿到需要还钱的 tid，from，money
	rows,err := tx.Queryx(BackSql)
	if err != nil {
		log.Println("获取还钱字段失败：",err)
		tx.Rollback()
		return
	}
	defer rows.Close()

	//3.拿到所有的rows记录后，将记录映射到结构体上
	dest := []Model.TransferModel{}
	_ = sqlx.StructScan(rows,&dest)
	fmt.Println("dest=",dest)
	//拿到所有的还钱记录后，就可以根据 tid,from,money 通过sql语句更新账单表，将钱还上。

	for _,tm :=range dest {
		//还钱：
		_,err := tx.Exec(replay,tm.Money,tm.From)
		if err != nil {
			log.Println("还钱更新用户帐单报错:",err )
			tx.Rollback() // 如果还钱失败，则回滚
		}
		//还钱后，将isback字段设置为1，代表处理过还钱业务了。
		_,err = tx.Exec(isback,tm.Tid)
		if err != nil {
			log.Println("还钱更新isback字段报错：",err)
			tx.Rollback()
		}
	}

	//最后提交事务
	defer commitTx(tx)
}

/*   自定义还钱事务锁
	为什么要设置这样一个锁？
	- 因为随着计划任务的监测，会不断的有符合超时还钱的请求，为了防止第一个还钱没执行结束完第二个任务就开始了
	- 所以需要设置一个 bool锁，只有一个事务提交之后，将锁释放才能执行下一个.
*/
//封装提交事务的函数
func commitTx(tx *sqlx.Tx) {
	err := tx.Commit()
	//如果提交有错误，并且有其它的回滚和提交错误
	if err != nil && err != sql.ErrTxDone{
		log.Println("tx,error:",err)
	}
	//提交事务后，将锁释放
	Backlock = false
}


//补偿机制三，转账重发
func TransResend(){

	//1.获取能够重发的转账数据
	rows,err := AppInit.GetDB().Queryx(Retrans)
	if err != nil {
		log.Println("转账重发：",err)
	}

	//2.将重发数据写入结构体
	dest := []Model.TransferModel{}
	err = sqlx.StructScan(rows,&dest)
	fmt.Println("重发时的数据：",dest)
	if err != nil {
		log.Println("转账重发2：",err)
	}else{
		MQ := Lib.NewMQ()
		//3.如果没有错误，则循环遍历每一条数据，并重新发送给MQ
		for _,tm := range dest{
			data,_ := json.Marshal(tm)
			err := MQ.SendMessage(Lib.Exchange_Trans,Lib.Router_Transkey,string(data))
			if err != nil {
				log.Println(err)
			}else{
				log.Println("重发MQ成功",err)
				TransConfirm()	//重发成功后将 status = 1
				return
			}
		}
	}
}

//如果第一次消息就发送成功了，那就直接将 status =1
func TransConfirm() error {
	_,err := AppInit.GetDB().Exec(ConfirmSql)
	if err != nil {
		return err
	}
	return nil
}