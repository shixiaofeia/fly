package mq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Conf struct {
	Address  string
	Port     string
	User     string
	Password string
}
type reCoon func() (*amqp.Connection, error)

var (
	DefaultConn *amqp.Connection
	ReCoon      reCoon
)

// Init
func Init(c Conf) (err error) {
	ReCoon = func() (*amqp.Connection, error) {
		coon, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
			c.User,
			c.Password,
			c.Address,
			c.Port))
		return coon, err
	}
	return
}

type AmqpC struct {
	channel *amqp.Channel
}

// NewChannel 获取新的连接
func NewChannel() (r *AmqpC, err error) {
	r = &AmqpC{}
	if DefaultConn == nil || DefaultConn.IsClosed() {
		DefaultConn, err = ReCoon()
		if err != nil {
			return nil, fmt.Errorf(" New mq coon err: %v", err)
		}
	}
	r.channel, err = DefaultConn.Channel()
	if err != nil {
		return nil, fmt.Errorf(" Get mq channel err: %v", err)
	}

	return r, nil
}

// ExchangeDeclare 创建交换机
func (a *AmqpC) ExchangeDeclare(name string, kind string) (err error) {
	return a.channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布订阅
func (a *AmqpC) Publish(exchange, key string, body []byte) (err error) {
	return a.channel.Publish(exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
}

// QueueDeclare 创建队列
func (a *AmqpC) QueueDeclare(name string) (err error) {
	_, err = a.channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueBind 绑定队列
func (a *AmqpC) QueueBind(name, key, exchange string) (err error) {
	return a.channel.QueueBind(name, key, exchange, false, nil)
}

// Consume 消费队列
func (a *AmqpC) Consume(queue, consumer string, handler func([]byte) error) (err error) {
	var msgList <-chan amqp.Delivery
	msgList, err = a.channel.Consume(queue, consumer, false, false, false, false, nil)
	if err != nil {
		return
	}
	for msg := range msgList {
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}
	return nil
}
