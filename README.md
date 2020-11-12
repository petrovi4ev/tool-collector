# pubsub-collector

## Usage

```go
    // create client Redis and other code
    
    // create a context with the collector shutdown deadline
    ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
    
    // create a collector
    c := collector.New(client, channelName)

    // run the collector
    c.Run(ctx)
    
    /* 
        the tested code is executed 
        ...
    */
    
    // waiting for the deadline
    <-ctx.Done()
    
    // we get messages in the form of a map: <redis channel name> -> <messages>
    messages := c.Messages()
    
    // check these messages
    
    // we clear the collector of messages from all channels
    
    c.Clean()
```