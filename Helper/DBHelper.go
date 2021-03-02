package Helper

import (
	"RabbitMQ/AppInit"
	"log"
)

//返回
func SetNotify(userid string, max_retry int) (r int) {
	sql := `insert into user_notify(user_id,updatetime)
	values(?,now()) ON DUPLICATE KEY UPDATE
	notifynum = IF(isdone=1,notifynum,notifynum+1),
	isdone = IF(notifynum>=?,1,0),
	updatetime = IF(isdone=1,updatetime,now());`

	//
	stmt, err := AppInit.GetDB().Prepare(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	//
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
