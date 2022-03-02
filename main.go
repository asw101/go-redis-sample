package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return errors.New("environment variable REDIS_ADDR is required")
	}
	addrs := []string{}
	for _, a := range strings.Split(addr, " ") {
		if strings.Contains(a, ":") {
			addrs = append(addrs, a)
		} else {
			addrs = append(addrs, fmt.Sprintf("%s:6379", a))
		}
	}
	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		return errors.New("environment variable REDIS_PASSWORD is required")
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: password,
		DB:       0, // use default DB
	})

	cmd := "consume"
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
		fmt.Printf("%s\n", cmd)
	}

	switch cmd {
	case "consume":
		return consume(rdb)
	case "produce":
		return produce(rdb)
	}
	return nil
}

func produce(rdb redis.UniversalClient) error {
	for {
		now := time.Now().Format("15:04:05")
		i, err := rdb.LPush(ctx, "now", now).Result()
		if err != nil {
			return err
		}
		fmt.Printf("%d\n", i)
		time.Sleep(2 * time.Second)
	}
}

func consume(rdb redis.UniversalClient) error {
	for {
		res, err := rdb.RPop(ctx, "now").Result()
		if err == redis.Nil {
			seconds := 2
			log.Printf("list empty. waiting %d seconds.\n", seconds)
			time.Sleep(time.Duration(2) * time.Second)
		} else if err != nil {
			return err
		} else {
			fmt.Printf("%s\n", res)
		}
	}
}

// original example
func example(rdb redis.UniversalClient) error {
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		return err
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		return err
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		return err
	} else {
		fmt.Println("key2", val2)
	}
	return nil
}
