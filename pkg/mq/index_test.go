package mq

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func nowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TestChannel_Publish(t *testing.T) {
	var (
		conf = Conf{
			User: "guest",
			Pwd:  "guest",
			Addr: "127.0.0.1",
			Port: "5672",
		}

		exchangeName = "user.register.direct"
		queueName    = "user.register.queue"
		keyName      = "user.register.event"
	)

	if err := Init(conf); err != nil {
		log.Fatalf(" mq init err: %v", err)
	}

	ch := NewChannel()
	if err := ch.ExchangeDeclare(exchangeName, "direct"); err != nil {
		log.Fatalf("create exchange err: %v", err)
	}
	if err := ch.QueueDeclare(queueName); err != nil {
		log.Fatalf("create queue err: %v", err)
	}
	if err := ch.QueueBind(queueName, keyName, exchangeName); err != nil {
		log.Fatalf("bind queue err: %v", err)
	}

	go func() {
		if err := NewConsumer(queueName, func(body []byte) error {
			fmt.Println("consume msg :" + string(body))
			return nil
		}); err != nil {
			log.Fatalf("consume err: %v", err)
		}
	}()

	go func() {
		for {
			if err := ch.Publish(exchangeName, keyName, []byte(nowTime())); err != nil {
				log.Fatalf("publish msg err: %v", err)
			}
			time.Sleep(time.Second)
		}

	}()

	time.Sleep(time.Minute)
	t.Log("end")
}

func TestChannel_PublishWithDelay(t *testing.T) {
	var (
		conf = Conf{
			User: "guest",
			Pwd:  "guest",
			Addr: "127.0.0.1",
			Port: "5672",
		}

		exchangeName   = "user.delay.direct"
		queueName      = "user.delay.queue"
		delayQueueName = "user.delay1.queue" // 延迟队列
		keyName        = "user.delay.event"
		delayKeyName   = "user.delay1.event" // 延迟key
	)

	if err := Init(conf); err != nil {
		log.Fatalf(" mq init err: %v", err)

	}

	ch := NewChannel()
	if err := ch.ExchangeDeclare(exchangeName, "direct"); err != nil {
		log.Fatalf("create exchange err: %v", err)
	}
	if err := ch.QueueDeclare(queueName); err != nil {
		log.Fatalf("create queue err: %v", err)
	}
	if err := ch.QueueDeclareWithDelay(delayQueueName, exchangeName, keyName); err != nil {
		log.Fatalf("create queue err: %v", err)
	}
	if err := ch.QueueBind(queueName, keyName, exchangeName); err != nil {
		log.Fatalf("bind queue err: %v", err)
	}
	if err := ch.QueueBind(delayQueueName, delayKeyName, exchangeName); err != nil {
		log.Fatalf("bind queue err: %v", err)
	}

	go func() {
		if err := NewConsumer(queueName, func(body []byte) error {
			fmt.Println(fmt.Sprintf("consumer msg: %s, ts: %s", string(body), nowTime()))
			return nil
		}); err != nil {
			log.Fatalf("consume err: %v", err)
		}
	}()
	go func() {
		for {
			if err := ch.PublishWithDelay(exchangeName, delayKeyName, []byte(nowTime()), 10*time.Second); err != nil {
				log.Fatalf("publish msg err: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Minute)
	t.Log("end")
}
