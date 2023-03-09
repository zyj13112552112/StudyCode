package main

import (
	amqp "github.com/streadway/amqp"
	"log"
)

var (
	conn *amqp.Connection
	err  error
)

func main() {
	conn, err = amqp.Dial("amqp://admin:admin@42.194.222.25:5672/")
	if err != nil {
		log.Fatal("连接mq失败 ", err)
		return
	}
	defer conn.Close()
	receive()
}

func receive() {
	ch, _ := conn.Channel()
	defer ch.Close()
	//队列配置
	q, _ := ch.QueueDeclare("hello", false, false, false, false, nil)
	//消费信息,返回的msgs是一个通道
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)
	forever := make(chan bool)

	//从通道中获取
	go func() {
		for msg := range msgs {
			log.Println("接收到消息: ", string(msg.Body))
		}
	}()

	//等待
	<-forever

}
