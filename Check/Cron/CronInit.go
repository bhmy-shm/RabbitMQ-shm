package Cron

import (
	"RabbitMQ/Check/CheckA"
	"github.com/robfig/cron/v3"
)

var MyCron *cron.Cron

func InitCron() error{
	MyCron = cron.New(cron.WithSeconds())

	//1、超时过期
	_,err := MyCron.AddFunc("0/3 * * * * *",CheckA.TransTimeOut)
	if err != nil {
		return err
	}

	//2、超时还钱
	_,err2 := MyCron.AddFunc("0/4 * * * * *",CheckA.TransBackMoney)
	if err2 != nil {
		return err2
	}

	//3、转账重发
	_,err3 := MyCron.AddFunc("0/10 * * * * *",CheckA.TransResend)
	if err3 != nil {
		return err3
	}

	return nil
}