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
	send()
}

func send() {
	//通道
	ch, _ := conn.Channel()
	defer ch.Close()
	//队列配置
	q, _ := ch.QueueDeclare("hello", false, false, false, false, nil)
	//消息发送
	msg := "hello world"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "test/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		log.Fatal("消息发送失败 ", err)
		return
	}
	log.Println("消息发送成功，消息内容: ", msg)
}
