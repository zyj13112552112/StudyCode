package main

import (
	amqp "github.com/streadway/amqp"
	"log"
)

var (
	conn *amqp.Connection
	err  error
)

// 消费者需要先启动监听，然后再启动生产者，否则消息因为没有监听而被安全的丢弃
func main() {
	conn, err = amqp.Dial("amqp://admin:admin@42.194.222.25:5672/")
	if err != nil {
		log.Fatal("连接mq失败 ", err)
		return
	}
	defer conn.Close()
	//FanoutReceive()
	//DirectReceive()
	TopicReceive()
}

// 需要将队列绑定交换机
// fanout类型，绑定到交换机上的所有队列都可以收到消息
func FanoutReceive() {
	ch, _ := conn.Channel()
	defer ch.Close()
	//交换机声明
	_ = ch.ExchangeDeclare("fanoutExchange", "fanout", true, false, false, false, nil)
	//队列声明1
	q1, _ := ch.QueueDeclare("fantoutQue", false, false, false, false, nil)
	//交换机队列绑定1
	_ = ch.QueueBind(q1.Name, "", "fanoutExchange", false, nil)
	//消费1
	msgs1, _ := ch.Consume(q1.Name, "", true, false, false, false, nil)
	//从通道中获取1
	go func() {
		for msg := range msgs1 {
			log.Println("队列1接收到消息: ", string(msg.Body))
		}
	}()
	//队列声明2
	q2, _ := ch.QueueDeclare("fanoutQue2", false, false, false, false, nil)
	//队列绑定2
	_ = ch.QueueBind(q2.Name, "", "fanoutExchange", false, nil)
	//消费2
	msgs2, _ := ch.Consume(q2.Name, "", true, false, false, false, nil)
	go func() {
		for msg := range msgs2 {
			log.Println("队列2接收到消息: ", string(msg.Body))
		}
	}()

	forever := make(chan bool)
	//等待
	<-forever

}

// direct,利用key绑定队列与交换机，发送消息指定key即可发送到某个队列上
func DirectReceive() {
	ch, _ := conn.Channel()
	defer ch.Close()

	//交换机
	_ = ch.ExchangeDeclare("directExchange", "direct", true, false, false, false, nil)
	//队列1
	q1, _ := ch.QueueDeclare("que1", false, false, true, false, nil)
	//队列2
	q2, _ := ch.QueueDeclare("que2", false, false, true, false, nil)
	//队列3
	q3, _ := ch.QueueDeclare("que3", false, false, true, false, nil)
	//绑定
	_ = ch.QueueBind(q1.Name, "q1", "directExchange", false, nil)
	_ = ch.QueueBind(q2.Name, "q2", "directExchange", false, nil)
	_ = ch.QueueBind(q3.Name, "q3", "directExchange", false, nil)

	msgs1, _ := ch.Consume(q1.Name, "", true, false, false, false, nil)
	msgs2, _ := ch.Consume(q2.Name, "", true, false, false, false, nil)
	msgs3, _ := ch.Consume(q3.Name, "", true, false, false, false, nil)

	go func() {
		for msg := range msgs1 {
			log.Println("队列p1接收到消息:", string(msg.Body))
		}
	}()

	go func() {
		for msg := range msgs2 {
			log.Println("队列p2接收到消息:", string(msg.Body))
		}
	}()

	go func() {
		for msg := range msgs3 {
			log.Println("队列p3接收到消息:", string(msg.Body))
		}
	}()

	forever := make(chan bool)
	<-forever
}

// topic灵活路由匹配
func TopicReceive() {
	ch, _ := conn.Channel()
	defer ch.Close()

	//交换机
	_ = ch.ExchangeDeclare("topicExchange", "topic", true, false, false, false, nil)
	//队列1
	q1, _ := ch.QueueDeclare("", true, false, false, false, nil)
	//队列2
	q2, _ := ch.QueueDeclare("", true, false, false, false, nil)
	//队列3
	q3, _ := ch.QueueDeclare("", true, false, false, false, nil)
	//绑定
	_ = ch.QueueBind(q1.Name, "topic.*.p1", "topicExchange", false, nil)
	_ = ch.QueueBind(q2.Name, "topic.*.p2", "topicExchange", false, nil)
	_ = ch.QueueBind(q3.Name, "topic.*.p3", "topicExchange", false, nil)

	msgs1, err := ch.Consume(q1.Name, "", true, false, false, false, nil)
	msgs2, err := ch.Consume(q2.Name, "", true, false, false, false, nil)
	msgs3, err := ch.Consume(q3.Name, "", true, false, false, false, nil)

	if err != nil {
		log.Println(err)
	}
	go func() {
		for msg := range msgs1 {
			log.Println("队列p1接收到消息:", string(msg.Body))
		}
	}()

	go func() {
		for msg := range msgs2 {
			log.Println("队列p2接收到消息:", string(msg.Body))
		}
	}()

	go func() {
		for msg := range msgs3 {
			log.Println("队列p3接收到消息:", string(msg.Body))
		}
	}()

	forever := make(chan bool)
	<-forever
}
