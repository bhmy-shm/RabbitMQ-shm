package Lib

import "fmt"

func UserInit() error {
	MQ := NewMQ()
	if MQ == nil {
		return fmt.Errorf("mq init is nil ")
	}
	defer MQ.Channel.Close()

	//1.创建交换器
	err := MQ.Channel.ExchangeDeclare(Exchange, "direct", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Exchange Error :%v", err)
	}

	//2.根据创建成功的交换机生成队列
	queues := fmt.Sprintf("%v,%v", QueueName, QueueName_union)
	err = MQ.DecQueueAndBind(queues, RouteKey, Exchange)
	if err != nil {
		return fmt.Errorf("Queue Bind error:%v", err)
	}
	return nil
}
