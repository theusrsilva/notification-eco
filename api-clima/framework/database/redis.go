package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(db int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Erro ao conectar no Redis: %v", err)
	}
	fmt.Println("Conectado ao Redis!")
	return &Redis{Client: client, Ctx: ctx}
}
func (r *Redis) Insert(key string, value string, expiration time.Duration) error {
	err := r.Client.Set(r.Ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("Erro ao definir chave no Redis: %w", err)
	}
	return nil
}
func (r *Redis) Find(key string) (string, error) {
	val, err := r.Client.Get(r.Ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("Erro ao buscar chave no Redis: %w", err)
	}
	return val, nil
}
