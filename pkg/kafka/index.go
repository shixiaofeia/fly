package kafka

import (
	"context"
	"fmt"
	kafkago "github.com/segmentio/kafka-go"
	"time"
)

type Config struct {
	Addr []string
}

var cfg *Config

func Init(c Config) error {
	cfg = &c
	if cfg == nil || len(cfg.Addr) == 0 {
		return fmt.Errorf("cfg is nil")
	}
	conn, err := kafkago.Dial("tcp", cfg.Addr[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// CreateTopic 创建主题.
func CreateTopic(topic string, partition int) (*kafkago.Conn, error) {
	if cfg == nil || len(cfg.Addr) == 0 {
		return nil, fmt.Errorf("cfg is nil")
	}
	conn, err := kafkago.DialLeader(context.Background(), "tcp", cfg.Addr[0], topic, partition)
	if err != nil {
		return nil, fmt.Errorf("dial leader err: %v", err)
	}

	return conn, nil
}

// NewWriterAsync 异步写入, 高性能, WriteMessages不会阻塞, 错误异步处理.
func NewWriterAsync(topic string, Completion func([]kafkago.Message, error)) *kafkago.Writer {
	w := &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Addr...),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		Async:        true,
		Completion:   Completion,
		RequiredAcks: kafkago.RequireOne,
	}

	return w
}

// NewWriterSync 同步写入, 性能较低, WriteMessages会阻塞直至处理完成, 能实时返回错误.
func NewWriterSync(topic string) *kafkago.Writer {
	w := &kafkago.Writer{
		Addr:     kafkago.TCP(cfg.Addr...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
		Async:    false,
	}

	return w
}

// NewReaderAutoCommit 消费自动确认.
func NewReaderAutoCommit(ctx context.Context, topic, groupId string, handle func(kafkago.Message) error) error {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Addr,
		Topic:    topic,
		MaxWait:  time.Second,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read msg err: %v", err)
		}
		if err = handle(msg); err != nil {
			return err
		}
	}
}

// NewReaderAckCommit 消费成功确认.
func NewReaderAckCommit(ctx context.Context, topic, groupId string, handle func(kafkago.Message) error) error {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Addr,
		Topic:    topic,
		MaxWait:  time.Second,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch msg err: %v", err)
		}
		if err = handle(msg); err != nil {
			return err
		}
		if err = r.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit msg err: %v", err)
		}
	}
}
