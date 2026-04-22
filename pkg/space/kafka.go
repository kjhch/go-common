package space

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaRegistrant interface {
	Handlers() map[string]KafkaTopicHandler
}
type KafkaTopicHandler = func(ctx context.Context, m kafka.Message) error

func NewKafkaWriter(cl *ConfigLoader) *kafka.Writer {
	if cl.injectConf.Mq.Addr == "" {
		return nil
	}
	return &kafka.Writer{
		Addr:                   kafka.TCP(cl.injectConf.Mq.Addr),
		Async:                  true, // 异步
		AllowAutoTopicCreation: true,
		BatchTimeout:           100 * time.Millisecond,
	}
}

type KafkaListener struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	cl         *ConfigLoader
	logger     *Logger
	registrant KafkaRegistrant
}

func NewKafkaListener(cl *ConfigLoader, logger *Logger, registrant KafkaRegistrant) *KafkaListener {
	if cl.injectConf.Mq.Addr == "" {
		return nil
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &KafkaListener{
		ctx:        ctx,
		cancel:     cancelFunc,
		cl:         cl,
		logger:     logger,
		registrant: registrant,
	}
}

func (kl *KafkaListener) Start() {
	for topic, handler := range kl.registrant.Handlers() {
		kl.wg.Go(func() {
			kl.handleTopic(topic, handler)
		})
		kl.logger.Info("Kafka消费者已启动, 主题: " + topic)
	}
	kl.wg.Wait()
}

func (kl *KafkaListener) handleTopic(topic string, handler KafkaTopicHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kl.cl.injectConf.Mq.Addr},
		Topic:   topic,
		GroupID: kl.cl.injectConf.Mq.GroupId, // 指定消费者组id
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
		ctx := context.Background()
		for _, h := range msg.Headers {
			ctx = context.WithValue(ctx, h.Key, string(h.Value))
		}

		if slices.Contains(kl.cl.injectConf.Log.Mq.EnabledTopics, msg.Topic) {
			kl.logger.InfoContext(ctx, fmt.Sprintf("[kafka]收到消息, topic:%s, msg:%s", msg.Topic, msg.Value))
		}

		err = handler(ctx, msg)

		if err != nil {
			kl.logger.ErrorContext(ctx, "[kafka]消息处理失败", "err", err)
		}
	}
	err := r.Close()
	if err != nil {
		kl.logger.Error("Kafka消费者异常退出, 主题: "+topic, "err", err)
		return
	}
	kl.logger.Info("Kafka消费者已关闭, 主题: " + topic)
}

func (kl *KafkaListener) Stop() {
	kl.logger.Info("Kafka消费者关闭中...")
	kl.cancel()
}
