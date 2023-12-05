package kafka

import (
	"context"
	"fmt"
	kafkago "github.com/segmentio/kafka-go"
	"log"
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

type (
	// consumerFunc 消费者func.
	consumerFunc  func(ctx context.Context, param consumerParam) error
	consumerParam struct {
		topic   string                      // topic 主题名称.
		groupID string                      // groupID 消费者组ID.
		handle  func(kafkago.Message) error // handle 处理消息.
	}
)

// reconnect 自动重连.
func reconnect(ctx context.Context, param consumerParam, f consumerFunc) error {
	num := 1
	for {
		if err := f(ctx, param); err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Printf("NewReaderAutoCommit topic: %s kafka err: %v, reconnect reader num: %d", param.topic, err.Error(), num)
				num += 1
				time.Sleep(3 * time.Second)
			}
		}
	}
}

// NewReaderAutoCommit 消费自动确认.
func NewReaderAutoCommit(ctx context.Context, topic, groupID string, handle func(kafkago.Message) error) error {
	return reconnect(ctx, consumerParam{
		topic:   topic,
		groupID: groupID,
		handle:  handle,
	}, readerAutoCommit)
}

// readerAutoCommit 消费自动确认.
func readerAutoCommit(ctx context.Context, param consumerParam) error {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Addr,
		Topic:    param.topic,
		MaxWait:  time.Second,
		GroupID:  param.groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read msg err: %v", err)
		}
		if err = param.handle(msg); err != nil {
			return err
		}
	}
}

// NewReaderAckCommit 消费成功确认.
func NewReaderAckCommit(ctx context.Context, topic, groupID string, handle func(kafkago.Message) error) error {
	return reconnect(ctx, consumerParam{
		topic:   topic,
		groupID: groupID,
		handle:  handle,
	}, readerAckCommit)
}

// readerAckCommit 消费成功确认.
func readerAckCommit(ctx context.Context, param consumerParam) error {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Addr,
		Topic:    param.topic,
		MaxWait:  time.Second,
		GroupID:  param.groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch msg err: %v", err)
		}
		if err = param.handle(msg); err != nil {
			return err
		}
		if err = r.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit msg err: %v", err)
		}
	}
}
