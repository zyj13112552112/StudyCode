package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func main() {
	dial, _ := amqp.Dial("amqp://admin:admin@42.194.222.25:5672/")
	defer dial.Close()

	ch, _ := dial.Channel()
	defer ch.Close()

	//死信交换机
	ch.ExchangeDeclare("dlx", "direct", false,
		true, false, false, nil)

	//正常交换机
	ch.ExchangeDeclare("normal", "fanout", false,
		true, false, false, nil)

	//队列参数
	var args = make(map[string]any)
	//设置队列过期时间
	args["x-message-ttl"] = 5000
	//设置死信交换机
	args["x-dead-letter-exchange"] = "dlx"
	//设置死信交换机key
	args["x-dead-letter-routing-key"] = "dlxkey"

	//正常队列
	normalque, _ := ch.QueueDeclare("normal", false, true,
		false, true, args) //这里传入死信配置参数args
	ch.QueueBind(normalque.Name, "normalkey", "normal", false, nil)

	//死信队列
	dlxque, _ := ch.QueueDeclare("dlx", false,
		true, false, true, nil)
	ch.QueueBind(dlxque.Name, "dlxkey", "dlx", false, nil)

	msgs, _ := ch.Consume(dlxque.Name, "", true, false, false, false, nil)
	go func() {
		for msg := range msgs {
			fmt.Println(string(msg.Body))
		}
	}()

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("test dlx message .%d", i)
		ch.Publish("normal", "normalkey", false,
			false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		time.Sleep(2 * time.Second)
	}

	forever := make(chan bool)
	<-forever

}
