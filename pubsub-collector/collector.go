package pubsub_collector

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

type Collector interface {
	Run(ctx context.Context)
	Messages() ChannelMessageMap
	Clean()
}

type PubSubMsgCollector struct {
	redisClient  *redis.Client
	redisChannel string
	messages     PubSubMessages
}

func New(redisClient *redis.Client, redisChannel string) Collector {
	messages := PubSubMessages{M: make(ChannelMessageMap), mx: &sync.Mutex{}}

	return &PubSubMsgCollector{
		redisClient:  redisClient,
		redisChannel: redisChannel,
		messages:     messages,
	}
}

func (collector *PubSubMsgCollector) Run(ctx context.Context) {
	go func(ctx context.Context) {
		pubsub, cancel := collector.subscribe()
		defer cancel()

		go collector.collect(pubsub)

		<-ctx.Done()
	}(ctx)
}

func (collector *PubSubMsgCollector) Clean() {
	collector.messages.mx.Lock()
	defer collector.messages.mx.Unlock()

	collector.messages.M = make(ChannelMessageMap, 0)
}

func (collector *PubSubMsgCollector) subscribe() (pubsub *redis.PubSub, cancel func()) {
	pubsub = collector.redisClient.Subscribe(context.Background(), collector.redisChannel)

	cancel = func() {
		if err := pubsub.Close(); err != nil {
			fmt.Printf("subscribe closing error: %s", err.Error())
		}
	}

	return
}

func (collector *PubSubMsgCollector) collect(pubsub *redis.PubSub) {
	for {
		msgi, err := pubsub.Receive(context.Background())

		if err != nil {
			break
		}

		select {
		default:
			switch msg := msgi.(type) {
			case *redis.Subscription:
			case *redis.Message:
				collector.messages.Store(msg.Channel, msg.Payload)
			default:
				fmt.Printf("error: unknown message: %#v", msgi)
			}
		}
	}
}
