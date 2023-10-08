package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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
	_, err := client.SAdd(ctx, "users_"+chatroomPath, displayName).Result()
	if err != nil {
		log.Println("Error adding user to chatroom:", err)
	}
	return err
}

func isUserInChatroomRedis(ctx context.Context, client *redis.Client, displayname string, chatroomPath string) (bool, error) {
	isMember, err := client.SIsMember(ctx, "users_"+chatroomPath, displayname).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func removeUserFromChatroomRedis(ctx context.Context, client *redis.Client, displayName string, chatroomPath string) error {
	_, err := client.SRem(ctx, "users_"+chatroomPath, displayName).Result()
	if err != nil {
		log.Println("Error removing user from chatroom:", err)
	}
	return err
}

func changeUserNameRedis(ctx context.Context, client *redis.Client, oldName string, newName string, chatroomPath string) error {

	userExists, err := isUserInChatroomRedis(ctx, client, oldName, chatroomPath)

	if userExists && err == nil {
		err := removeUserFromChatroomRedis(ctx, client, oldName, chatroomPath)
		if err != nil {
			return err
		} else {
			err := addUserToChatroomRedis(ctx, client, newName, chatroomPath)
			if err != nil {
				return err
			}
		}
	} else if !userExists {
		log.Panicf("%v not found in Redis db", oldName)
	}

	return err

}

func getMessageHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string) ([]string, error) {
	chatMessages, err := client.LRange(ctx, "messages_"+chatroomPath, 0, -1).Result()
	if err != nil {
		log.Println("Error getting message history from chatroom:", err)
	}
	return chatMessages, err
}

func getMessageHistoryLengthRedis(ctx context.Context, client *redis.Client, chatroomPath string) (int64, error) {
	length, err := client.LLen(ctx, "messages_"+chatroomPath).Result()
	if err != nil {
		log.Println("Error getting length of message history from chatroom:", err)
	}
	return length, err
}

func addMessageToHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string, chatMessage string) error {
	_, err := client.RPush(ctx, "messages_"+chatroomPath, chatMessage).Result()
	if err != nil {
		log.Println("Error adding message to chatroom history:", err)
	}
	return err
}

func removeMessageFromHistoryRedis(ctx context.Context, client *redis.Client, chatroomPath string) error {
	_, err := client.LPop(ctx, "messages_"+chatroomPath).Result()
	if err != nil {
		log.Println("Error removing message from chatroom history:", err)
	}
	return err
}

func deleteKeyRedis(ctx context.Context, client *redis.Client, key string) error {
	_, err := client.Del(ctx, "messages_"+key).Result()
	if err != nil {
		log.Println("Error deleting message history of chatroom:", err)
	}
	return err
}

func getAllChatroomsRedis(ctx context.Context, client *redis.Client, key string) (map[string]string, error) {
	rooms, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Println("Error getting all chatrooms in lobby:", err)
	}
	return rooms, err
}

func addChatroomToLobbyRedis(ctx context.Context, client *redis.Client, key string, room map[string]interface{}) error {
	_, err := client.HSet(ctx, key, room).Result()
	if err != nil {
		log.Println("Error adding chatroom to lobby:", err)
	}
	return err
}

func removeChatroomFromLobbyRedis(ctx context.Context, client *redis.Client, key string, roomname string) error {
	_, err := client.HDel(ctx, key, roomname).Result()
	if err != nil {
		log.Println("Error removing chatroom from lobby:", err)
	}
	return err
}

func connectCockrochDB(context context.Context) *pgxpool.Pool {

	conn, err := pgxpool.New(context, os.Getenv("COCKROACHDB"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connection to CockroachDB established")
	}
	return conn

}

func createTableIfNotExistsCRDB(ctx context.Context, client *pgxpool.Pool, tblName string, tblSchema string) {
	_, err := client.Query(ctx, "CREATE TABLE IF NOT EXISTS "+tblName+" "+tblSchema)
	if err != nil {
		log.Println(err)
	}
}

func addUserToChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayName string, chatroomPath string) error {

	createTableIfNotExistsCRDB(ctx, client, "users", "(chatroompath STRING, displayname STRING)")

	_, err := client.Query(ctx, "INSERT INTO users (chatroompath, displayname) VALUES ($1, $2)", chatroomPath, displayName)
	if err != nil {
		log.Println("Error adding user to chatroom:", err)
	}
	return err

}

// func isUserInChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayname string, chatroomPath string) (bool, error) {
// 	isMember, err := client.SIsMember(ctx, "users_"+chatroomPath, displayname).Result()
// 	if err != nil {
// 		return false, err
// 	}

// 	return isMember, nil
// }

// func removeUserFromChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayName string, chatroomPath string) error {
// 	_, err := client.SRem(ctx, "users_"+chatroomPath, displayName).Result()
// 	if err != nil {
// 		log.Println("Error removing user from chatroom:", err)
// 	}
// 	return err
// }

// func changeUserNameCRDB(ctx context.Context, client *pgxpool.Pool, oldName string, newName string, chatroomPath string) error {

// 	userExists, err := isUserInChatroomRedis(ctx, client, oldName, chatroomPath)

// 	if userExists && err == nil {
// 		err := removeUserFromChatroomRedis(ctx, client, oldName, chatroomPath)
// 		if err != nil {
// 			return err
// 		} else {
// 			err := addUserToChatroomRedis(ctx, client, newName, chatroomPath)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	} else if !userExists {
// 		log.Panicf("%v not found in Redis db", oldName)
// 	}

// 	return err

// }

// func getMessageHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) ([]string, error) {
// 	chatMessages, err := client.LRange(ctx, "messages_"+chatroomPath, 0, -1).Result()
// 	if err != nil {
// 		log.Println("Error getting message history from chatroom:", err)
// 	}
// 	return chatMessages, err
// }

// func getMessageHistoryLengthCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) (int64, error) {
// 	length, err := client.LLen(ctx, "messages_"+chatroomPath).Result()
// 	if err != nil {
// 		log.Println("Error getting length of message history from chatroom:", err)
// 	}
// 	return length, err
// }

// func addMessageToHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string, chatMessage string) error {
// 	_, err := client.RPush(ctx, "messages_"+chatroomPath, chatMessage).Result()
// 	if err != nil {
// 		log.Println("Error adding message to chatroom history:", err)
// 	}
// 	return err
// }

// func removeMessageFromHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) error {
// 	_, err := client.LPop(ctx, "messages_"+chatroomPath).Result()
// 	if err != nil {
// 		log.Println("Error removing message from chatroom history:", err)
// 	}
// 	return err
// }

// func deleteKeyCRDB(ctx context.Context, client *pgxpool.Pool, key string) error {
// 	_, err := client.Del(ctx, "messages_"+key).Result()
// 	if err != nil {
// 		log.Println("Error deleting message history of chatroom:", err)
// 	}
// 	return err
// }

// func getAllChatroomsCRDB(ctx context.Context, client *pgxpool.Pool, key string) (map[string]string, error) {
// 	rooms, err := client.HGetAll(ctx, key).Result()
// 	if err != nil {
// 		log.Println("Error getting all chatrooms in lobby:", err)
// 	}
// 	return rooms, err
// }

// func addChatroomToLobbyCRDB(ctx context.Context, client *pgxpool.Pool, key string, room map[string]interface{}) error {
// 	_, err := client.HSet(ctx, key, room).Result()
// 	if err != nil {
// 		log.Println("Error adding chatroom to lobby:", err)
// 	}
// 	return err
// }

// func removeChatroomFromLobbyCRDB(ctx context.Context, client *pgxpool.Pool, key string, roomname string) error {
// 	_, err := client.HDel(ctx, key, roomname).Result()
// 	if err != nil {
// 		log.Println("Error removing chatroom from lobby:", err)
// 	}
// 	return err
// }
