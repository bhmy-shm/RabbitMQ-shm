package main

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Check/CheckA"
	"RabbitMQ/Helper"
	"RabbitMQ/Lib"
	"RabbitMQ/Middle"
	"RabbitMQ/Model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

//A公司的服务端启动


func main(){

	r := gin.Default()
	r.Use(Middle.ErrorRecover())	//引用异常捕获中间件

	//转账路由
	r.Handle("POST", "/", func(context *gin.Context) {

		//1.初始化专门用来获取POST数据的结构体
		TransferModel := Model.NewTransfer()
		err := context.BindJSON(&TransferModel)
		//2.如果没有接收到信息则返回错误
		Middle.CheckError(err,"BindJSON获取参数失败")

		//3.拿到提交消息后进行转账
		err = Helper.TransMoney(TransferModel)
		Middle.CheckError(err,"转账失败")

		//4.如果接收到消息则将转账信息写入到消息队列
		MQ := Lib.NewMQ()
		data,_:=json.Marshal(TransferModel)
		//fmt.Println("写入消息队列的数据：",string(data))
		err = MQ.SendMessage(Lib.Exchange_Trans,Lib.Router_Transkey,string(data))
		Middle.CheckError(err,"写入消息队列失败")
		//5.这里写入发送消息到队列后，就将status变成1，否则后期会重发冗余消息队列
		err = CheckA.TransConfirm()
		if err != nil {
			log.Println("TransConfirm err :",err)
		}

		//6.返回给web浏览器一个页面
		context.JSON(200,gin.H{"result":TransferModel})
	})

	//回调路由
	r.Handle("POST","/callback", func(context *gin.Context) {
		tid := context.PostForm("tid")
		sql := "update a.translog set isback=1 where tid=? and status=1"
		ret,err := AppInit.GetDB().Exec(sql,tid)
		if err != nil {
			log.Println("r.Handle GetDB.Exec() --- err:",err)
		}

		affCount,err2 := ret.RowsAffected()
		fmt.Println("callback affcount=",int(affCount))
		if err2 != nil || affCount != 1{
			context.String(200,"error")
		}else{
			context.String(200,"success")
		}
	})

	//二、go程开启允许服务
	c := make(chan error)
	go func(){
		err := r.Run(":8080")
		if err != nil {
			c<-err
		}
	}()
	//初始化数据库
	go func(){
		err := AppInit.DBInit("a")
		if err != nil {
			c<-err
		}
	}()
	//初始化转账队列，在服务运行的一开始就将转账队列初始化
	go func(){
		err := Lib.TransInit()
		if err != nil {
			c <- err
		}
	}()
	err := <-c
	log.Fatal(err)
}
