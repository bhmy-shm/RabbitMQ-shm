package main

import (
	"RabbitMQ/Lib"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

//消费者
var c string

func init() {
	flag.StringVar(&c, "c", "", "消费者名称")
	flag.Parse()
}
func main() {

	if c == "" {
		log.Fatal("需要指定consumer参数：-c [name]")
	}

	//1.直接从管道中消费数据
	MQ := Lib.NewMQ()

	//2.设置消费者限流
	err := MQ.Channel.Qos(2, 0, false)
	if err != nil {
		log.Fatal(err)
	}
	MQ.Counsumer(Lib.QueueName, c, SendMail)
	defer MQ.Channel.Close()
}

//消费者发送邮件
func SendMail(msgs <-chan amqp.Delivery, consumer string) {
	for msg := range msgs {
		fmt.Println(consumer, "收到消息 ", string(msg.Body))
		go Send(c, msg)
	}
}

func Send(consumer string, msg amqp.Delivery) error {
	time.Sleep(time.Second * 1)
	fmt.Printf("%s向UserID=%s的用户发送邮件\n", consumer, string(msg.Body))
	msg.Ack(false)
	return nil
}

/*=================================================*/
//**模拟C1出现问题，continue跳过则不会执行 msg.ACK确认收到
// if c == "c1" {
// 	msg.Reject(false) //不重新入列，会将分配到c1的消息丢弃
// 	msg.Reject(true)  //允许重新入列，会将消息分配给其它的 消费者进行消费
// 	continue
// }

// func Consumer() {

// 	//1.连接RabbitMQ
// 	conn := AppInit.GetConn()
// 	defer conn.Close()

// 	//2.生成消息通道
// 	c, err := conn.Channel()
// 	if err != nil {
// 		log.Fatal("C", err)
// 	}
// 	defer c.Close()

// 	//3.生成consumer读取消息通道的数据
// 	//参数1：队列名称，参数2：消费者名称
// 	msgs, err := c.Consume("rabbit", "consumer1", false, false, false, false, nil)
// 	if err != nil {
// 		log.Fatal("Consumer", err)
// 	}

// 	//4.遍历consumer读取到的数据，c.Consumer返回的是一个管道
// 	for msg := range msgs {
// 		fmt.Printf("从%v消息队列中读取数据:%s\n", msg.DeliveryTag, string(msg.Body))
// 	}
// }
