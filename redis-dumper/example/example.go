package main

import (
	"fmt"
	redisdumper "github.com/BitMedia-IO/tool-collector/redis-dumper"
	"gopkg.in/redis.v3"
	"os"
)

func main() {
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

	if err := client.FlushAll().Err(); err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}

	if err := redisdumper.JSONImportRedisDB(client, "./example.json"); err != nil {
		fmt.Printf("import error: %s", err.Error())
		os.Exit(1)
	}
}
