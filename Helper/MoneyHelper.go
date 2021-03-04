package Helper

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Model"
	"context"
	"log"
)

/*	此文件是用来做转账业务处理，通过 “Mysql事务” 完成整个转账的过程
- 转账过程1：更新A公司账户表，扣款。如果扣款失败，则Rollback回滚，成功则进入下一步。
- 转账过程2：向转账日志表写入转账的消息，通过 rows判断是否成功。
- 转账过程3：如果全部成功，则commit提交，完成事务。*/

var(
	//扣款的sql语句,  :money  :from   是占位符，配合NamedExec使用
	DeductMoney = "update usermoney set user_money=user_money-:money " +
		"where user_name=:from and user_money>=:money"

	//记录转账日志的sql语句，
	TransferLog = "insert into translog(`from`,`to`,money,updatetime) "+
		"values(:from,:to,:money,now())"
)

//转账
func TransMoney(tm *Model.TransferModel) error {

	//1.开启事务
	tx := AppInit.GetDB().MustBeginTx(context.Background(),nil)

	// ================  业务一：扣款 ==================

	//2.在事务中进行扣款
	//NamedExec事务中的命名查询。任何命名的占位符参数都将被来自arg的字段替换。
	ret,_:=tx.NamedExec(DeductMoney,tm)

	//3.如果执行扣款的sql语句没有问题，则提交。如果有问题，则回滚.
	rowAffected,_ := ret.RowsAffected()
	if rowAffected == 0{
		err := tx.Rollback() 	//回滚
		if err != nil { log.Fatal("扣款失败") ; return err }
	}

	// ================= 业务二：记录转账日志 =================

	//4.扣款结束后记录转账日志
	ret,_ = tx.NamedExec(TransferLog,tm)
	rowAffected,_ = ret.RowsAffected()
	if rowAffected == 0{
		err := tx.Rollback() 	//回滚
		if err != nil { log.Fatal("插入日志记录失败") ; return err }
	}

	// ================ 业务三：=========================
	//5.记录日志之后，将交易号tid赋值给B
	tid,_ := ret.LastInsertId()	//赋值tid到B公司
	tm.Tid = int(tid)

	return tx.Commit()
}

