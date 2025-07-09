package cmd

import (
	"context"
	"log"
	"marketflow/config"
	"marketflow/internal/adapter/redis"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"time"
)

func FillRedisWithData() {
	ctx := context.Background()
	cfg := config.Redis{
		Addr:     "localhost:6379",
		Password: "strongpassword",
	}
	client, err := redis.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	testData := domain.PriceData{
		Symbol:    types.ETHUSDT,
		Exchange:  types.Exchange1,
		Price:     99.99,
		Timestamp: time.Now(),
	}

	if err := client.SetLatest(ctx, testData, time.Minute*5); err != nil {
		log.Fatal(err)
	}

	log.Println("set data", testData)
}
