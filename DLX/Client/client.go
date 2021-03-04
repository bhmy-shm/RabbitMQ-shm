package main

import (
	"RabbitMQ/DLX"
	"RabbitMQ/Lib"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//消费者
func TestDlx(msgs <-chan amqp.Delivery,c string){
	for msg:=range msgs{
		fmt.Println("收到消息",string(msg.Body))
		msg.Ack(false)
	}
}

func main()  {
	client:=Lib.NewMQ()
	err:=DLX.TestInit()
	if err!=nil{
		log.Fatal(err)
	}
	err=client.Channel.Qos(2,0,false)
	if err!=nil{
		log.Fatal(err)
	}
	//消费订阅测试队列
	client.Counsumer("test_queue","c",TestDlx)
	defer client.Channel.Close()
}