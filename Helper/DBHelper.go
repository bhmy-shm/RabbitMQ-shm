package Helper

import (
	"RabbitMQ/AppInit"
	"log"
)

//该文件是延迟队列的帮助文件，不断的向数据库中插入数据，通过 notifynum | isdone 字段进行判断
//如果 notifynum = 5，则不会再继续插入数据，最终的 rows结果就会是0。


//批量插入，根据返回的rows进行判断
func SetNotify(userid string, max_retry int) (r int) {
	sql := `insert into user_notify(user_id,updatetime)
	values(?,now()) ON DUPLICATE KEY UPDATE
	notifynum = IF(isdone=1,notifynum,notifynum+1),
	isdone = IF(notifynum>=?,1,0),
	updatetime = IF(isdone=1,updatetime,now());`

	//1、数据库预处理
	stmt, err := AppInit.GetDB().Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	//2、传入注册信息的userid，和notifynum的最大值，配合预处理的sql语句，进行插入数据。
	ret, err := stmt.Exec(userid, max_retry)
	if err != nil {
		log.Println("fetch notify error :", err)
		return
	}
	//受影响的行 >0 继续发送，==0代表到了最大次数，不再发送
	Affected, err := ret.RowsAffected()
	if err != nil {
		log.Println("RowsAffected error:", err)
		return
	}
	return int(Affected)
}
