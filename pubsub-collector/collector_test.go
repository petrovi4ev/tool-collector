package pubsub_collector_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSingleMessage(t *testing.T) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	channelName := "test"
	testMsg := fmt.Sprintf("Ok, guys. Now is %d", time.Now().UnixNano())

	t.Run("not empty message", func(t *testing.T) {
		testCollector.Run(ctx)
		time.Sleep(1 * time.Second)

		redisClient.Publish(ctx, channelName, testMsg)

		<-ctx.Done()

		assert.Equal(t, 1, len(testCollector.Messages()))
		assert.Equal(t, 1, len(testCollector.Messages()[channelName]))
		assert.Equal(t, testMsg, testCollector.Messages()[channelName][0])
	})
	t.Run("empty message", func(t *testing.T) {
		testCollector.Clean()
		testCollector.Run(ctx)
		redisClient.Publish(context.Background(), channelName, "")

		<-ctx.Done()

		assert.Equal(t, 0, len(testCollector.Messages()))
	})
}
