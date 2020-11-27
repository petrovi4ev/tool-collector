package main

import (
	"context"
	"fmt"
	collector "github.com/BitMedia-IO/tool-collector/pubsub-collector"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
	"time"
)

func main() {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	channelName := "test"

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

	c := collector.New(client, channelName)
	c.Run(ctx)
	time.Sleep(50 * time.Millisecond)

	go func(ctx context.Context) {
		tick := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-tick.C:
				nanoTime := time.Now().UnixNano()
				fmt.Println(nanoTime)
				client.Publish(context.Background(), channelName, fmt.Sprintf("Ok, guys. Now is %d", nanoTime))
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	<-ctx.Done()

	messages := c.Messages()
	fmt.Printf("Messages:\n\t%+v", strings.Join(messages[channelName], "\n\t")+"\n")

	fmt.Println("Clean messages...")
	c.Clean()
	fmt.Println("Messages: ", c.Messages())
}
