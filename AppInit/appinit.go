package AppInit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var MQConn *amqp.Connection

type MQUser struct {
	Name string
	Pass string
	Host string
	Port uint
}

func init() {
	muser := MQUser{Name: "shm", Pass: "123.com", Host: "192.168.168.4", Port: 5672}
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d", muser.Name, muser.Pass, muser.Host, muser.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}
	MQConn = conn
	log.Print(MQConn.Major) //服务器的主要版本
}

func GetConn() *amqp.Connection {
	return MQConn
}
