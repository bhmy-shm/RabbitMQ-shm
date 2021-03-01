package Lib

import (
	"RabbitMQ/AppInit"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

const (
	QueueName       = "user"         //消息队列1，用户注册
	QueueName_union = "user_union"   //消息队列2，其它合作单位的用户注册
	Exchange        = "UserExChange" //用户注册的交换器
	RouteKey        = "userreg"      //交换器路由键
)

type MQ struct {
	Channel       *amqp.Channel          //生产者消费者的消息通道
	notifyConfirm chan amqp.Confirmation //Confirm模式的通道
	notifyReturn  chan amqp.Return       //NotifyReturn模式的通道
}

func NewMQ() *MQ {
	conn := AppInit.GetConn()
	c, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return &MQ{Channel: c}
}

//生产者
func (this *MQ) SendMessage(exchange string, key string, message string) error {
	//4.发送队列到交换器
	err := this.Channel.Publish(exchange, key, true, false, amqp.Publishing{
		ContentType: "text/palin",
		Body:        []byte(message),
	})
	return err
}

//消费者
func (this *MQ) Counsumer(queue string, key string, callback func(<-chan amqp.Delivery, string)) error {
	msgs, err := this.Channel.Consume(queue, key, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	callback(msgs, key)
	return nil
}

//循环生成队列
func (this *MQ) DecQueueAndBind(queues string, key string, exchange string) error {
	//分隔队列名称
	qList := strings.Split(queues, ",")
	//1.循环创建多个队列
	for _, queue := range qList {
		q, err := this.Channel.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			return err
		}

		//2.每创建一个队列就绑定一个路由键
		err = this.Channel.QueueBind(q.Name, key, exchange, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

//开启服务端的 Confirm 模式
func (this *MQ) SetConfirm() {
	err := this.Channel.Confirm(false)
	if err != nil {
		log.Println(err)
	}
	this.notifyConfirm = this.Channel.NotifyPublish(make(chan amqp.Confirmation))

}

//接收Confirm管道返回的信息
func (this *MQ) ListenConfirm() {
	ret := <-this.notifyConfirm
	defer this.Channel.Close()
	if ret.Ack {
		log.Println("Confire消息发送成功", ret)
	} else {
		log.Println("Confirm消息发送失败", ret)
	}
}

func (this *MQ) NotifyReturn() {
	//初始化Return的channel
	this.notifyReturn = this.Channel.NotifyReturn(make(chan amqp.Return))
	//如果消息没有进入正确的队列，则会向 Return 返回只一个值
	go this.ListenReturn() //使用协程执行
}

func (this *MQ) ListenReturn() {
	ret := <-this.notifyReturn
	fmt.Println("111", string(ret.Body))
	if string(ret.Body) != "" {
		fmt.Println("到这")
		log.Println("当前消息没有正确入列", string(ret.Body))
	}
}
