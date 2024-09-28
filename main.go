package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	tlsEnabled := os.Getenv("REDIS_TLS_ENABLED")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	clientCacheEnabled := os.Getenv("REDIS_CLIENT_CACHE_ENABLED")

	var tlsConfig *tls.Config
	if tlsEnabled == "1" || tlsEnabled == "true" {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client, err := rueidisotel.NewClient(rueidis.ClientOption{
		TLSConfig:    tlsConfig,
		Username:     username,
		Password:     password,
		InitAddress:  []string{net.JoinHostPort(redisHost, redisPort)},
		SelectDB:     int(0),
		DisableCache: !(clientCacheEnabled == "1" || clientCacheEnabled == "true"),
	})
	if err != nil {
		fmt.Println("error creating client:", err)
		os.Exit(1)
	}

	ctx, cancel1 := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel1()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("db ping successful")
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel2()

	if err := client.Do(ctx, client.B().Set().Key("key").Value("val").Build()).Error(); err != nil {
		fmt.Println("error setting key:", err)
		os.Exit(1)
	}

	ctx, cancel3 := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel3()

	v, err := client.Do(ctx, client.B().Get().Key("key").Build()).ToString()
	if err != nil {
		fmt.Println("error getting key:", err)
		os.Exit(1)
	}

	fmt.Println(v)
}
