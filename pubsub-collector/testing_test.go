package pubsub_collector_test

import (
	"fmt"
	collector "github.com/BitMedia-IO/tool-collector/pubsub-collector"
	"gopkg.in/redis.v3"
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
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	return collector.New(redisClient, "test")
}

func getTestRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6384",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Printf("connection error: %s", err.Error())
		os.Exit(1)
	}

	return client
}
