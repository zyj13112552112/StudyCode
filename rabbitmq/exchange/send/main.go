package main

import (
	"fmt"
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
	//Fantoutsend()
	//Directsend()
	TopicSend()
}

// fanout类似广播，只要监听该交换机的消费者都可以收到消息
func Fantoutsend() {
	//通道
	ch, _ := conn.Channel()
	defer ch.Close()
	//交换机配置
	_ = ch.ExchangeDeclare("fanoutExchange", "fanout", true, false, false, false, nil)
	err = ch.Publish("fanoutExchange", "", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("这是一条广播消息"),
		})
	if err != nil {
		log.Fatal("消息发送失败，", err)
	} else {
		log.Println("消息发送成功")
	}
}

// direct,利用key绑定队列与交换机，发送消息指定key即可发送到某个队列上
func Directsend() {
	ch, _ := conn.Channel()
	defer ch.Close()

	//交换机配置
	_ = ch.ExchangeDeclare("directExchange", "direct", true, false, false, false, nil)

	for i := 0; i < 5; i++ {
		ch.Publish("directExchange", "q1", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
	}

	for i := 0; i < 5; i++ {
		ch.Publish("directExchange", "q2", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
	}

	for i := 0; i < 5; i++ {
		ch.Publish("directExchange", "q3", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
	}
}

// topic灵活路由匹配
func TopicSend() {
	ch, _ := conn.Channel()
	defer ch.Close()

	//交换机配置
	_ = ch.ExchangeDeclare("topicExchange", "topic", true, false, false, false, nil)

	for i := 0; i < 5; i++ {
		err = ch.Publish("topicExchange", "topic.get.p1", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
		if err != nil {
			log.Println(err)
		}
	}

	for i := 0; i < 5; i++ {
		err = ch.Publish("topicExchange", "topic.get.p2", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
		if err != nil {
			log.Println(err)
		}
	}

	for i := 0; i < 5; i++ {
		err = ch.Publish("topicExchange", "topic.get.p3", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("这是消息%d", i)),
		})
		if err != nil {
			log.Println(err)
		}
	}
}
