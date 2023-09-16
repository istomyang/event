package sender

import (
	"git.grizzlychina.com/infrastructures/event"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

type kafka struct {
	config        KafkaConfig
	producerAsync sarama.AsyncProducer
	producerSync  sarama.SyncProducer
	messagesBatch []*sarama.ProducerMessage
	preSendTime   time.Time
	mut           sync.Mutex
}

type KafkaConfig struct {
	Hosts                  []string       `json:"hosts" yaml:"hosts"`
	Topic                  string         `json:"topic" yaml:"topic"`
	MessagesPerSend        int            `json:"messages_per_send,omitempty" yaml:"messages_per_send,omitempty"`
	MessagesPerSendTimeout time.Duration  `json:"messages_per_send_timeout,omitempty" yaml:"messages_per_send_timeout,omitempty"`
	Extra                  *sarama.Config `json:"extra" yaml:"extra"`
	Encode                 func(m event.Message) []byte
	ErrHandler             func(error)
}

var _ event.Sender = &kafka{}

func NewKafkaSender(config KafkaConfig) event.Sender {
	return &kafka{
		config: config,
	}
}

func (k *kafka) Send(message event.Message) {
	var sendMessage = &sarama.ProducerMessage{
		Topic: k.config.Topic,
		Value: sarama.ByteEncoder(k.config.Encode(message)),
	}
lo:
	for {
		select {
		case k.producerAsync.Input() <- sendMessage:
			break lo
		case err := <-k.producerAsync.Errors():
			k.config.ErrHandler(err)
		}
	}
}

func (k *kafka) SendSync(message event.Message) error {
	var sendMessage = &sarama.ProducerMessage{
		Topic: k.config.Topic,
		Value: sarama.ByteEncoder(k.config.Encode(message)),
	}
	var err error

	if k.config.MessagesPerSend <= 1 {
		_, _, err = k.producerSync.SendMessage(sendMessage)
		return err
	}

	k.mut.Lock()
	if len(k.messagesBatch) > k.config.MessagesPerSend || time.Now().After(k.preSendTime) {
		var m = k.messagesBatch
		m = append(m, sendMessage)
		k.messagesBatch = nil
		k.preSendTime = time.Now()
		k.mut.Unlock()
		err = k.producerSync.SendMessages(m)
		return err
	}
	k.messagesBatch = append(k.messagesBatch, sendMessage)
	k.mut.Unlock()
	return err
}

func (k *kafka) Run() error {
	{
		producer, err := sarama.NewAsyncProducer(k.config.Hosts, k.config.Extra)
		if err != nil {
			return err
		}
		k.producerAsync = producer
	}
	{
		producer, err := sarama.NewSyncProducer(k.config.Hosts, k.config.Extra)
		if err != nil {
			return err
		}
		k.producerSync = producer
	}

	return nil
}

func (k *kafka) MustRun() {
	if err := k.Run(); err != nil {
		panic(err)
	}
}

func (k *kafka) Close() error {
	var err error
	err = k.producerAsync.Close()
	err = k.producerSync.Close()
	return err
}
