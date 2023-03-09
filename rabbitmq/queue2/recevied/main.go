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
	//设置预取计数值为1。告诉RabbitMQ一次只向一个worker发送一条消息。换句话说，在处理并确认前一个消息之前，不要向工作人员发送新消息。
	//设置当消息未被消费，不向work发送消息
	ch.Qos(1, 0, false)
	//队列配置
	q, _ := ch.QueueDeclare("hello", false, false, false, false, nil)
	//消费信息,返回的msgs是一个通道,关闭自动确认
	msgs1, _ := ch.Consume(q.Name, "", false, false, false, false, nil)

	msgs2, _ := ch.Consume(q.Name, "", false, false, false, false, nil)

	forever := make(chan bool)

	//从通道中获取
	go func() {
		for msg := range msgs1 {
			log.Println("1接收到消息: ", string(msg.Body))
			msg.Ack(false) //自动确认
		}
	}()

	go func() {
		for msg := range msgs2 {
			log.Println("2接收到消息: ", string(msg.Body))
			msg.Ack(false) //自动确认
		}
	}()
	//
	//go func() {
	//	for msg := range msgs2 {
	//		log.Println("2接收到消息: ", string(msg.Body))
	//		//msg.Ack(false) //自动确认
	//	}
	//}()

	//等待
	<-forever

}
