package main

import (
	"RabbitMQ/DLX"
	"RabbitMQ/Lib"
	"github.com/gin-gonic/gin"
	"log"
)

func main(){

	router:=gin.Default()
	router.Handle("POST","/dxl", func(context *gin.Context) {
		//发送到MQ
		mq:=Lib.NewMQ()
		err:=mq.SendMessage(DLX.ExName,"topic.#","delay message")
		if err!=nil{
			log.Println(err)
		}
		context.JSON(200,gin.H{"result":"延迟消息发送成功"})
	})
	c:=make(chan error)
	go func() {
		err:=router.Run(":8080")
		if err!=nil{
			c<-err
		}
	}()
	go func() {
		err:=DLX.DXLInit()
		if err!=nil{
			c<-err
		}
	}()
	err:=<-c
	log.Fatal(err)
}
