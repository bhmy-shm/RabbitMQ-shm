package main

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Check/Cron"
	"log"
)

/*   开启补偿机制，补偿机制需要通过 Cron来进行监测
	监测1：转账超时
*/

func main(){

	c := make(chan error)

	//初始化数据库
	go func() {
		err := AppInit.DBInit("a")
		if err != nil { c <- err}
	}()

	//初始化定时任务，开启
	go func(){
		err := Cron.InitCron()
		if err != nil { c <- err}
		Cron.MyCron.Start() //开启Cron计划任务
	}()
	err := <-c
	log.Fatal(err)
}

