package CheckB

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//CheckB,B公司的订阅补偿

var (
	//1.将消费到的数据记录到数据库的转账日志表当中
	SaveSql = "insert into translog(tid,`from`,`to`,money,updatetime) values(?,?,?,?,now())"

	//2.收到转账日志记录后，确认收钱。
	Proceeds = "update umoney set user_money=user_money + ? where user_name = ?"

	//3.回调地址
	CallbackUrl = "http://127.0.0.1:8080/callback"
)

//一、B公司获取转账数据
func SaveLog(tm *Model.TransferModel){

	//1.开启事务
	tx,err := AppInit.GetDB().BeginTxx(context.Background(),nil)
	if err != nil {
		log.Println("B tx error:",err)
		return
	}

	//2.在事务中将获取到的数据写入B公司数据库日志表
	_,err = AppInit.GetDB().Exec(SaveSql,tm.Tid,tm.From,tm.To,tm.Money)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}else{
		log.Println("写入B公司数据库日志表成功")
	}

	//3.如果写入日志表成功，则根据日志表的字段，进行收款(需要传入钱数，和收款人)
	ret,err := tx.Exec(Proceeds,tm.Money,tm.To)
	if err != nil {
		log.Println("Proceeds Exec err:",err)
		return
	}else{
		log.Println("收款成功")
	}

	aff,err_aff := ret.RowsAffected()  //rowaffected返回受更新、插入或删除影响的行数。
	if err_aff != nil || aff != 1 {
		log.Println("到这erraff",err_aff)
		tx.Rollback()
		return
	}else{
		log.Println("获取受影响的行数:",int(aff))
	}

	//4.收款成功后进行回调，将A公司的 status = 结果修改成1
	fmt.Println("tm.tid=",tm.Tid)
	err = callBack(tm.Tid)
	if err != nil {
		log.Println("Callback error:",err)
		tx.Rollback()
		return
	}else{
		log.Println("回调成功")
	}

	//5.提交事务
	defer commitTx(tx)
}

//二、B公司消费者订阅队列
func SubFromA(msgs <-chan amqp.Delivery,c string){

	//1.从Delivery管道中获取队列数据
	for msg := range msgs{
		tm := Model.NewTransfer()
		err := json.Unmarshal(msg.Body,tm) //结果：zhangsan转给lisi,金额:20
		if err != nil {
			log.Println(err)
		}else{
			//2.开启go程，将收到的数据保存到数据库，并完成确认交付。
			go func(t *Model.TransferModel){
				SaveLog(t)
				defer msg.Ack(false)
			}(tm)
		}
	}
}

//三、收款成功后的回调函数
func callBack(tid int) error{
	rsp,err:=http.Post("http://127.0.0.1:8080/callback",
		"application/x-www-form-urlencoded",strings.NewReader(fmt.Sprintf("tid=%d",tid)))
	if err != nil {
		log.Println("http.Post failed err =",err)
		return err
	}
	defer rsp.Body.Close()

	b,err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll failed err:",err)
		return err
	}
	fmt.Println("callBack() -- rsp.Body --- b =",string(b))

	if string(b) == "success" {
		return nil
	}else{
		return fmt.Errorf("error")
	}
}
//四、收款事务的提交封装
func commitTx(tx *sqlx.Tx){ //清理事务
	err:=tx.Commit()
	if err!=nil && err!=sql.ErrTxDone{
		log.Println("tx error:",err)
	}
}