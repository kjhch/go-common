package space

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"sync"
)

type KafkaRegistrant interface {
	Handlers() map[string]func(m kafka.Message)
	Brokers() []string
	GroupId() string
}
type KafkaListener struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	logger     *Logger
	registrant KafkaRegistrant
}

func NewKafkaListener(logger *Logger, registrant KafkaRegistrant) *KafkaListener {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &KafkaListener{
		ctx:        ctx,
		cancel:     cancelFunc,
		logger:     logger,
		registrant: registrant,
	}
}

func (kl *KafkaListener) Start() {
	for topic, handler := range kl.registrant.Handlers() {
		kl.wg.Go(func() {
			kl.handleTopic(topic, handler)
		})
		kl.logger.Info("Kafka消费者已启动", "topic", topic)
	}
	kl.wg.Wait()
}

func (kl *KafkaListener) handleTopic(topic string, handler func(m kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kl.registrant.Brokers(),
		Topic:   topic,
		GroupID: kl.registrant.GroupId(), // 指定消费者组id
		//RebalanceTimeout: time.Second,
		//MaxBytes: 10e6, // 10MB
	})

	for {
		msg, err := r.ReadMessage(kl.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			panic(err)
		}
		handler(msg)
	}
	kl.logger.Info("Kafka消费者已关闭", "topic", topic)
}

func (kl *KafkaListener) Stop() {
	kl.logger.Info("Kafka消费者关闭中...")
	kl.cancel()
}
