package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Conf struct {
		Address  string
		Port     string
		User     string
		Password string
	}
)

var (
	defaultConn    *Connection
	defaultChannel *Channel
)

// Init 初始化
func Init(c Conf) (err error) {
	defaultConn, err = Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.User,
		c.Password,
		c.Address,
		c.Port))
	if err != nil {
		return fmt.Errorf("new mq conn err: %v", err)
	}

	defaultChannel, err = defaultConn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %v", err)
	}
	return
}

// NewChannel 获取channel.
func NewChannel() *Channel {
	return defaultChannel
}

// ExchangeDeclare 创建交换机.
func (ch *Channel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布订阅.
func (ch *Channel) Publish(exchange, key string, body []byte) (err error) {
	return ch.Channel.Publish(exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
}

// QueueDeclare 创建队列.
func (ch *Channel) QueueDeclare(name string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueBind 绑定队列.
func (ch *Channel) QueueBind(name, key, exchange string) (err error) {
	return ch.Channel.QueueBind(name, key, exchange, false, nil)
}

// NewConsumer 实例化一个消费者, 会单独用一个channel.
func NewConsumer(queue string, handler func([]byte) error) error {
	ch, err := defaultConn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %v", err)
	}

	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, queue)
	}

	for msg := range deliveries {
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}
