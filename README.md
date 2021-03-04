# RabbitMQ-shm
RabbitMQ学习代码

### Resend-Message分支
___
利用延迟队列，实现消费者订阅消息后发送邮件的 “延迟重发”。


创建延时队列的代码：Lib / UserInit.go

关键的SQL语句：Helper / DBHelper.go

```sql
func SetNotify(userid string, max_retry int) (r int) {
    sql := `insert into user_notify(user_id,updatetime)
    values(?,now()) ON DUPLICATE KEY UPDATE
    notifynum = IF(isdone=1,notifynum,notifynum+1),
    isdone = IF(notifynum>=?,1,0),
    updatetime = IF(isdone=1,updatetime,now());`
```

实现重发业务逻辑的代码：Cmd / Client / Client.go




[博客详解](http://www.bhmy.top/blog/?p=1584)

- 延迟队列使用
  - [安装延迟队列](http://www.bhmy.top/blog/?p=1584#安装延迟队列)
  - [创建延迟队列](http://www.bhmy.top/blog/?p=1584#创建延迟队列)
- 利用延迟队列，实现邮件延迟重发
  - [关键点：SQL语句模拟事务](http://www.bhmy.top/blog/?p=1584#关键点：SQL语句模拟事务)
  - [邮件失败后重发](http://www.bhmy.top/blog/?p=1584#邮件失败后重发)
  - [数据库设计](http://www.bhmy.top/blog/?p=1584#数据库设计)
