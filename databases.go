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

	log.Println("PING ->")
	state, err := client.Ping(context).Result()
	log.Println(state)
	if err != nil {
		log.Fatal(err)
	} else if state == "PONG" {
		log.Println("connection to Redis established")
	} else {
		log.Println("problem encountered connecting to Redis")
	}

	return client

}

func addUserToChatroomRedis(ctx context.Context, client *redis.Client, displayName string, chatroomPath string) error {
	_, err := client.SAdd(ctx, "chatroom:"+chatroomPath, displayName).Result()
	if err != nil {
		log.Println("Error adding user to chatroom:", err)
	}
	return err
}

// func getChatroomUsersRedis(ctx context.Context, client *redis.Client, chatroomPath string) ([]string, error) {
// 	members, err := client.SMembers(ctx, "chatroom:"+chatroomPath).Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return members, nil
// }

func isUserInChatroomRedis(ctx context.Context, client *redis.Client, displayname string, chatroomPath string) (bool, error) {
	isMember, err := client.SIsMember(ctx, "chatroom:"+chatroomPath, displayname).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func removeUserFromChatroomRedis(ctx context.Context, client *redis.Client, displayName string, chatroomPath string) error {
	_, err := client.SRem(ctx, "chatroom:"+chatroomPath, displayName).Result()
	if err != nil {
		log.Println("Error removing user from chatroom:", err)
	}
	return err
}

func changeUserNameRedis(ctx context.Context, client *redis.Client, oldName string, newName string, chatroomPath string) error {

	userExists, err := isUserInChatroomRedis(ctx, client, oldName, chatroomPath)

	if userExists && err != nil {
		err := removeUserFromChatroomRedis(ctx, client, oldName, chatroomPath)
		if err != nil {
			return err
		} else {
			err := addUserToChatroomRedis(ctx, client, newName, chatroomPath)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// func getMessageHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string) []string {
// 	chatMessages, err := client.LRange(ctx, chatroomPath, 0, -1).Result()
// 	if err != nil {
// 		log.Println("Error getting message history from chatroom:", err)
// 	}
// 	return chatMessages
// }

// func getMessageHistoryLengthRedis(ctx context.Context, client *redis.Client, chatroomPath string) int64 {
// 	length, err := client.LLen(ctx, chatroomPath).Result()
// 	if err != nil {
// 		log.Println("Error getting length of message history from chatroom:", err)
// 	}
// 	return length
// }

// func addMessageToHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string, chatMessage string) error {
// 	//create a message record interface/struct/json obj
// 	_, err := client.RPush(ctx, chatroomPath, chatMessage).Result()
// 	return err
// }

// func removeMessageFromHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string) error {
// 	_, err := client.LPop(ctx, chatroomPath).Result()
// 	return err
// }

// func deleteKeyRedis(ctx context.Context, client *redis.Client, key string) error {
// 	_, err := client.Del(ctx, key).Result()
// 	return err
// }

func connectCockrochDB(context context.Context) *pgx.Conn {

	conn, err := pgx.Connect(context, os.Getenv("COCKROACHDB"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection to CockroachDB established")
	}
	return conn

}
