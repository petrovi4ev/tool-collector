package pubsub_collector_test

import (
	"context"
	"fmt"
	collector "github.com/BitMedia-IO/tool-collector/pubsub-collector"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
)

var (
	testCollector collector.Collector
	redisClient   *redis.Client
)

func TestMain(m *testing.M) {
	testCollector = getTestCollector()
	redisClient = getTestRedisClient()

	m.Run()
}

func getTestCollector() collector.Collector {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6384"})
	pong, err := redisClient.Ping(context.Background()).Result()
	fmt.Println(pong, err)

	return collector.New(redisClient, "test")
}

func getTestRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6384",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("connection error: %s", err.Error())
		os.Exit(1)
	}

	return client
}
