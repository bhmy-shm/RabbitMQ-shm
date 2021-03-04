package main

import (
	"RabbitMQ/AppInit"
	"RabbitMQ/Check/CheckB"
	"RabbitMQ/Lib"
	"flag"
	"log"
)

var (
	c string
)

//消费者参数，指定消费者名称
func init(){
	flag.StringVar(&c,"c","","消费者名称")
	flag.Parse()
	if c == ""{
		log.Panic("请指定消费者名称：-c Name")
	}
}

func main(){
	err := AppInit.DBInit("b")
	if err != nil {
		log.Fatal("B公司，初始化DBInit失败：",err)
	}

	//1、创建一个QOS订阅限制
	MyClient := Lib.NewMQ()
	err = MyClient.Channel.Qos(2,0,false)
	if err != nil {
		log.Fatal("订阅限制err：",err)
	}

	//2、开启订阅
	MyClient.Counsumer(Lib.Queue_Trans,c,CheckB.SubFromA)
	defer MyClient.Channel.Close()
}