# RabbitMQ-shm
RabbitMQ学习代码

### Master分支

___

跨库的转账业务，通过补偿机制，以及sql + 事务的方式完成转账业务碰到的各种问题。



#### [RabbitMQ – 异构系统转账分布式事务](http://www.bhmy.top/blog/?p=1593)

- [数据库设计](http://www.bhmy.top/blog/?p=1593#数据库设计)
- [实现转账业务](http://www.bhmy.top/blog/?p=1593#实现转账业务)
  - [A公司转账业务逻辑，记录交易日志](http://www.bhmy.top/blog/?p=1593#A公司转账业务逻辑，记录交易日志)
  - [A公司转账业务逻辑，记录日之后发送消息MQ](http://www.bhmy.top/blog/?p=1593#A公司转账业务逻辑，记录日之后发送消息MQ)



#### [RabbitMQ – 转账业务的补偿机制](http://www.bhmy.top/blog/?p=1594)

- [补偿机制（一）转账超时失败](http://www.bhmy.top/blog/?p=1594#补偿机制（一）转账超时失败)
- [补偿机制（二），交易超时失败后还钱](http://www.bhmy.top/blog/?p=1594#补偿机制（二），交易超时失败后还钱)
- [补偿机制（三）延迟重发](http://www.bhmy.top/blog/?p=1594#补偿机制（三）延迟重发)
- [B公司开始消费](http://www.bhmy.top/blog/?p=1594#B公司开始消费)



### Resend-Message分支

___

利用延迟队列，实现消费者订阅消息后发送邮件的 “延迟重发”。



#### [RabbitMQ – 延迟队列](http://www.bhmy.top/blog/?p=1584)

- 延迟队列使用
  - [安装延迟队列](http://www.bhmy.top/blog/?p=1584#安装延迟队列)
  - [创建延迟队列](http://www.bhmy.top/blog/?p=1584#创建延迟队列)
- 利用延迟队列，实现邮件延迟重发
  - [关键点：SQL语句模拟事务](http://www.bhmy.top/blog/?p=1584#关键点：SQL语句模拟事务)
  - [邮件失败后重发](http://www.bhmy.top/blog/?p=1584#邮件失败后重发)
  - [数据库设计](http://www.bhmy.top/blog/?p=1584#数据库设计)



### Registration 分支

___

模拟用户注册流程。



#### [RabbitMQ – Gin 模拟用户注册流程写入MQ](http://www.bhmy.top/blog/?p=1543)

- [简单API过程、注册流程、MQ操作简单封装](http://www.bhmy.top/blog/?p=1543#简单API过程、注册流程、MQ操作简单封装)
- [整理生产者入队代码](http://www.bhmy.top/blog/?p=1543#整理生产者入队代码)
- [客户端消费（注册用户信息，确认消息）](http://www.bhmy.top/blog/?p=1543#客户端消费（注册用户信息，确认消息）)
- [多消费者消费消息，Reject()重新入列](http://www.bhmy.top/blog/?p=1543#多消费者消费消息，Reject()重新入列)
- [消费者限流](http://www.bhmy.top/blog/?p=1543#消费者限流)
- [开启Confirm模式，记录发送消息是否成功](http://www.bhmy.top/blog/?p=1543#开启Confirm模式，记录发送消息是否成功)
- [NotifyReturn 模式](http://www.bhmy.top/blog/?p=1543#NotifyReturn_模式)



#### 最后致谢沈逸老师...

