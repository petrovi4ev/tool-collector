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
    // since the collector is started in the goroutine, you need to wait a bit for it to start working
    time.Sleep(20 * time.Millisecond)
    
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

# mongo migrator docker [mongo\_migrations](./mongo_migrations)

https://github.com/seppevs/migrate-mongo

```
pushd mongo_migrations
docker image build -f Dockerfile.migrator -t eth_migrate:1.0 .

docker container run --network="host" eth_migrate:1.0 'mongodb://127.0.0.1:27017/Leads?authSource=admin'
docker container run eth_migrate:1.0  'mongodb://host.docker.internal:27017/Leads?authSource=admin'
```
