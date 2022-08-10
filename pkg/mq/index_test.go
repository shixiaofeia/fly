package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	var (
		conf = Conf{
			User:     "guest",
			Password: "guest",
			Address:  "127.0.0.1",
			Port:     "5672",
		}

		exchangeName = "user.register.direct"
		queueName    = "user.register.queue"
		keyName      = "user.register.event"
	)

	if err := Init(conf); err != nil {
		t.Errorf(" mq init err: %v", err)
		return
	}

	ch := NewChannel()
	if err := ch.ExchangeDeclare(exchangeName, "direct"); err != nil {
		t.Errorf("create exchange err: %v", err)
		return
	}
	if err := ch.QueueDeclare(queueName); err != nil {
		t.Errorf("create queue err: %v", err)
		return
	}
	if err := ch.QueueBind(queueName, keyName, exchangeName); err != nil {
		t.Errorf("bind queue err: %v", err)
	}

	go func() {
		if err := ch.Consume(queueName, "", func(body []byte) error {
			fmt.Println("consume msg :" + string(body))
			return nil
		}); err != nil {
			t.Errorf("consume err: %v", err)
			return
		}
	}()

	if err := ch.Publish(exchangeName, keyName, []byte("hello word")); err != nil {
		t.Errorf("publish msg err: %v", err)
	}

	time.Sleep(time.Minute)
	t.Log("success")
}
