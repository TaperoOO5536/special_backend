package kafka

import (
	"errors"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	errUnknownType = errors.New("unkwond event type")
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic, headerStr string) error {
	headers := []kafka.Header{
		{Key: "event-type", Value: []byte(headerStr)},
	}
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Headers: headers,
		Key: nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}
	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default: 
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(5000)
	p.producer.Close()
}