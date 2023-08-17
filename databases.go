package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

func connectRedis(context context.Context) *redis.Client {

	redisAddr := os.Getenv("REDISADDRESS")
	redisPass := os.Getenv("REDISPASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	state, err := client.Ping(context).Result()
	if err != nil {
		log.Fatal(err)
	} else if state == "PONG" {
		log.Println("connection to Redis established")
	} else {
		log.Printf("PING -> %v", state)
	}

	return client

}

func addUserToChatroomRedis(ctx context.Context, client *redis.Client, displayName, chatroomPath string) error {
	_, err := client.SAdd(ctx, "chatroom:"+chatroomPath, displayName).Result()
	if err != nil {
		log.Println("Error adding user to chatroom:", err)
	}
	return err
}

func connectCockrochDB(context context.Context) *pgx.Conn {

	conn, err := pgx.Connect(context, os.Getenv("COCKROACHDB"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection to CockroachDB established")
	}
	return conn

}
