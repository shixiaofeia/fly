package kafka

import (
	"context"
	"fmt"
	kafkago "github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func conf() Config {
	return Config{Addr: []string{"127.0.0.1:9092"}}
}

var (
	topic = "fly"
)

func TestCreateTopic(t *testing.T) {
	_ = Init(conf())
	if _, err := CreateTopic(topic, 0); err != nil {
		t.Error(err.Error())
	}
}

func TestNewWriterAsync(t *testing.T) {
	_ = Init(conf())
	w := NewWriterAsync(topic, func(messages []kafkago.Message, err error) {
		fmt.Println("async err", err, len(messages))
	})
	msg := make([]kafkago.Message, 0)
	for i := 0; i < 10; i++ {
		msg = append(msg, kafkago.Message{Value: []byte(fmt.Sprintf("num: %d", i))})
	}
	if err := w.WriteMessages(context.Background(), msg...); err != nil {
		t.Error(err.Error())
	}
	time.Sleep(2 * time.Second)

}
func TestNewWriterSync(t *testing.T) {
	_ = Init(conf())
	w := NewWriterSync(topic)
	msg := make([]kafkago.Message, 0)
	for i := 0; i < 5; i++ {
		msg = append(msg, kafkago.Message{Value: []byte(fmt.Sprintf("num: %d", i))})
	}
	if err := w.WriteMessages(context.Background(), msg...); err != nil {
		t.Error(err.Error())
	}
}

func TestNewReaderAutoCommit(t *testing.T) {
	_ = Init(conf())
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	if err := NewReaderAutoCommit(ctx, topic, "test1", func(message kafkago.Message) error {
		fmt.Println(string(message.Value))
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestNewReaderAckCommit(t *testing.T) {
	_ = Init(conf())
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	if err := NewReaderAckCommit(ctx, topic, "test2", func(message kafkago.Message) error {
		fmt.Println(string(message.Value))
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}
