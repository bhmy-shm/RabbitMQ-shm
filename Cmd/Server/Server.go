package main

import (
	"RabbitMQ/Lib"
	"RabbitMQ/Model"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//作为服务端，生产数据，也就是 producer

func main() {

	r := gin.New()
	r.Handle("POST", "/user", func(c *gin.Context) {

		//1.获取提交的用户注册信息
		usermodel := Model.NewUserModel()
		err := c.BindJSON(&usermodel)

		//2.判断是否接收到注册信息
		if err != nil {
			c.JSON(400, gin.H{"result": "param error"})
		} else { //如果收到消息则将用户ID写入消息队列
			usermodel.UserID = int(time.Now().Unix())
			if usermodel.UserID > 0 {
				mq := Lib.NewMQ() //连接RabbitMQ

				mq.SetConfirm()   //开启Confirm模式，监听是否能够发送消息
				mq.NotifyReturn() //监听Notifyreturn，监听消息是否能够入队

				err := mq.SendMessage(Lib.Exchange, Lib.RouteKey, strconv.Itoa(usermodel.UserID))
				if err != nil {
					log.Println("Server.go,error :=", err)
				}

				mq.ListenConfirm()
				// defer mq.Channel.Close()
				c.JSON(200, gin.H{"result": usermodel})
			}
		}
	})

	c := make(chan error)
	go func() {
		if err := Lib.UserInit(); err != nil {
			c <- err
		}
	}()
	go func() {
		r.Run(":8080")
	}()
	err := <-c
	log.Fatal(err)
}

// func producer() {
// 	//1.连接RabbitMQ
// 	conn := AppInit.GetConn()
// 	defer conn.Close()

// 	//2.生成消息通道
// 	c, err := conn.Channel()
// 	if err != nil {
// 		log.Fatal("C", err)
// 	}
// 	defer c.Close()

// 	//3.生成队列
// 	q, err := c.QueueDeclare("rabbit", false, false, false, false, nil)
// 	if err != nil {
// 		log.Fatal("Q", err)
// 	}

// 	//4.发送消息队列到MQ
// 	err = c.Publish("", q.Name, false, false, amqp.Publishing{
// 		ContentType: "text/plain",
// 		Body:        []byte("test01"),
// 	})
// 	if err != nil {
// 		log.Fatal("P", err)
// 	}
// 	log.Println("消息发送成功")
// }
