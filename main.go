package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/redis/rueidis"
)

func main() {
	redisURL, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("error parsing redis URL:", err)
		os.Exit(1)
	}

	redisHostPort := redisURL.Host
	redisUsername := redisURL.User.Username()
	redisPassword, _ := redisURL.User.Password()

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{redisHostPort},
		Username:    redisUsername,
		Password:    redisPassword,
	})
	if err != nil {
		fmt.Println("error connecting to redis:", err)
		os.Exit(1)
	}

	defer client.Close()

	ctx := context.Background()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("db ping successful")
	}

	if err := client.Do(ctx, client.B().Set().Key("key").Value("val").Build()).Error(); err != nil {
		fmt.Println("error setting key:", err)
		os.Exit(1)
	}

	v, err := client.Do(ctx, client.B().Get().Key("key").Build()).ToString()
	if err != nil {
		fmt.Println("error getting key:", err)
		os.Exit(1)
	}

	fmt.Println(v)
}
