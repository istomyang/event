package receiver

import (
	"git.grizzlychina.com/infrastructures/event"
	"github.com/IBM/sarama"
	"sync/atomic"
)

type kafka struct {
	config      *KafkaConfig
	consumer    sarama.Consumer
	messageChan chan event.Message
	chanClosed  atomic.Bool
}

var _ event.Receiver = &kafka{}

type KafkaConfig struct {
	Hosts  []string       `json:"hosts" yaml:"hosts"`
	Topic  string         `json:"topic" yaml:"topic"`
	Extra  *sarama.Config `json:"extra" yaml:"extra"`
	Encode func([]byte) event.Message
}

func NewKafkaReceiver(config KafkaConfig) event.Receiver {
	return &kafka{
		config:      &config,
		messageChan: make(chan event.Message),
	}
}

func (k *kafka) Receive() <-chan event.Message {
	return k.messageChan
}

func (k *kafka) Run() error {
	consumer, err := sarama.NewConsumer(k.config.Hosts, k.config.Extra)
	if err != nil {
		return err
	}

	k.consumer = consumer

	partition, err := k.consumer.Partitions(k.config.Topic)
	if err != nil {
		return err
	}

	for _, partitionID := range partition {
		consumer, err := k.consumer.ConsumePartition(k.config.Topic, partitionID, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		go func() {
			for message := range consumer.Messages() {
				if k.chanClosed.Load() {
					break
				}
				k.messageChan <- k.config.Encode(message.Value)
			}
		}()
	}
	return err
}

func (k *kafka) MustRun() {
	if err := k.Run(); err != nil {
		panic(err)
	}
}

func (k *kafka) Close() error {
	var err error
	err = k.consumer.Close()
	k.chanClosed.Store(true)
	close(k.messageChan)
	return err
}
