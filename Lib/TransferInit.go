package Lib

import "fmt"

func TransInit() error {
	//1.连接RabbitMQ，并生成消息通道
	MQ := NewMQ()
	if MQ == nil {
		return fmt.Errorf("mq init is nil ")
	}
	defer MQ.Channel.Close()

	//2.创建交换器
	err := MQ.Channel.ExchangeDeclare(Exchange_Trans, "direct", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Exchange Error :%v", err)
	}

	//3.根据创建成功的交换机生成队列
	//参数1:队列 ; 参数2：路由键 ; 参数3：交换机
	err = MQ.DecQueueAndBind(Queue_Trans, Router_Transkey, Exchange_Trans)
	if err != nil {
		return fmt.Errorf("Queue Bind error:%v", err)
	}
	return nil
}
