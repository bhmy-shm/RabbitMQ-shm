package DLX

import (
	"RabbitMQ/Lib"
	"fmt"
)

var(
	ExName = "Exchange_DLX"
	QueueName = "Queue_DLX"
)

//初始化死信队列
func DXLInit() error {
	MQ := Lib.NewMQ()
	defer MQ.Channel.Close()

	//一、申明交换器
	err := MQ.Channel.ExchangeDeclare(ExName,"direct",false,false,false,false,nil)
	if err != nil {
		return fmt.Errorf("Exchange_DLX error：",err)
	}

	//二、声明死信(中间)队列的固定参数
	args:=map[string]interface{}{"x-message-ttl":3000,	//TTL延迟时间，3000毫秒 = 3秒
		"x-dead-letter-exchange":"Exchange_test",	//死信指定的交换器，与测试的要相同
		"x-dead-letter-routing-key":"dlx"}	//死信指定路由键dlx，与测试的要相同

	//三、生成死信队列，传递args参数，绑定路由键，绑定交换器
	err = MQ.DecQueueAndBindWithArgs(QueueName,"topic.#",ExName,args)
	if err != nil {
		return fmt.Errorf("DLX Bind Error")
	}
	return nil
}

//初始化一个测试队列
func TestInit() error{
	MQ := Lib.NewMQ()
	defer MQ.Channel.Close()
	//申明交换器
	err := MQ.Channel.ExchangeDeclare("Exchange_test","direct",false,false,false,false,nil)
	if err != nil {
		return fmt.Errorf("test error",err)
	}
	//生成队列
	args := map[string]interface{}{}
	err = MQ.DecQueueAndBindWithArgs("test_queue","dlx","Exchange_test",args)
	if err != nil {
		return fmt.Errorf("dlx Bind error",err)
	}
	return nil
}


