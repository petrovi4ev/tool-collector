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
		time.Sleep(50 * time.Millisecond)

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

func TestManyMessages(t *testing.T) {
	channelName := "test"

	do := func(n int, emptyMsg bool) {
		for n > 0 {
			testMsg := fmt.Sprintf("Ok, guys. Now is %d", time.Now().UnixNano())
			if emptyMsg {
				testMsg = ""
			}
			redisClient.Publish(context.Background(), channelName, testMsg)
			time.Sleep(100 * time.Microsecond)

			n--
		}
	}

	t.Run("not empty messages", func(t *testing.T) {
		n := 100
		notEmptyCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))

		testCollector.Clean()
		testCollector.Run(notEmptyCtx)
		time.Sleep(50 * time.Millisecond)

		do(n, false)

		<-notEmptyCtx.Done()

		assert.Equal(t, 1, len(testCollector.Messages()))
		assert.Equal(t, n, len(testCollector.Messages()[channelName]))
	})
	t.Run("empty messages", func(t *testing.T) {
		n := 100
		emptyCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))

		testCollector.Clean()
		testCollector.Run(emptyCtx)
		time.Sleep(50 * time.Millisecond)

		do(n, true)

		<-emptyCtx.Done()

		assert.Equal(t, 1, len(testCollector.Messages()))
		assert.Equal(t, n, len(testCollector.Messages()[channelName]))
	})
	t.Run("mix messages", func(t *testing.T) {
		nEmpty := 100
		nNotEmpty := 100
		mixCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))

		testCollector.Clean()
		testCollector.Run(mixCtx)
		time.Sleep(50 * time.Millisecond)

		do(nNotEmpty, false)
		do(nEmpty, true)

		<-mixCtx.Done()

		assert.Equal(t, 1, len(testCollector.Messages()))
		assert.Equal(t, nNotEmpty+nEmpty, len(testCollector.Messages()[channelName]))
	})
}
