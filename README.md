# RabbitMQ-shm
RabbitMQ学习代码 - registration分支

模拟用户注册，实现注册成功后将用户ID发送给RabbitMQ，由消费者拿到用户ID后，对每一位用户ID发送邮件。

[详细文档地址](https://www.bhmy.top/blog/?p=1543)

- [简单API过程、注册流程、MQ操作简单封装](https://www.bhmy.top/blog/?p=1543#简单API过程、注册流程、MQ操作简单封装)
- [整理生产者入队代码](https://www.bhmy.top/blog/?p=1543#整理生产者入队代码)
- [客户端消费（注册用户信息，确认消息）](https://www.bhmy.top/blog/?p=1543#客户端消费（注册用户信息，确认消息）)
- [多消费者消费消息，Reject()重新入列](https://www.bhmy.top/blog/?p=1543#多消费者消费消息，Reject()重新入列)
- [消费者限流](https://www.bhmy.top/blog/?p=1543#消费者限流)
- [开启Confirm模式，记录发送消息是否成功](https://www.bhmy.top/blog/?p=1543#开启Confirm模式，记录发送消息是否成功)
- [NotifyReturn 模式](https://www.bhmy.top/blog/?p=1543#NotifyReturn_模式)


