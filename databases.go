package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectRedis(context context.Context, try int) *redis.Client {

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
		if try <= 5 {
			log.Printf("having trouble connecting to Redis, retrying(%v)...", try)
			client = connectRedis(context, try+1)
		}

	} else if state == "PONG" {
		log.Println("connection to Redis established")
	} else {
		log.Println("problem encountered connecting to Redis")
	}

	return client

}

func redisKeyExists(ctx context.Context, client *redis.Client, key string) bool {

	result, err := client.Exists(ctx, key).Result()
	if err != nil || result == 0 {
		log.Printf("%v Redis key does not exist. Error: %v", key, err)
		return false
	} else if result == 1 {
		log.Printf("%v Redis key exists", key)
		return true
	} else {
		log.Printf("%v Redis key does not exist. Unexpected result in query: %v", key, result)
		return false
	}

}

func resetRedisKeyExpiration(ctx context.Context, client *redis.Client, key string) {

	err := client.Expire(ctx, key, time.Duration(expirationSeconds)*time.Second).Err()
	if err != nil {
		log.Panicf("Error setting Redis Key expiration: %v", err)
	}

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

func getMessageHistoryLengthRedis(ctx context.Context, client *redis.Client, chatroomPath string) (int8, error) {
	length, err := client.LLen(ctx, "messages_"+chatroomPath).Result()
	if err != nil {
		log.Println("Error getting length of message history from chatroom:", err)
	}
	return int8(length), err
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
		log.Println("Error getting all chatrooms in lobby from Redis:", err)
	}
	return rooms, err
}

func addChatroomToLobbyRedis(ctx context.Context, client *redis.Client, key string, roomname string, roompath string) error {

	jsonroompath, err := json.Marshal(roompath)
	if err != nil {
		log.Panicf("Error marshaling JSON: %v", err)
		return err
	}

	room := map[string]interface{}{roomname: fmt.Sprintf("%v", string(jsonroompath))}

	_, err = client.HSet(ctx, key, room).Result()
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

	queryString := "CREATE TABLE IF NOT EXISTS " + tblName + " " + tblSchema

	_, err := client.Exec(ctx, queryString)
	if err != nil {
		log.Println(err)
		log.Printf("Query used: %v", queryString)
	}

}

func addUserToChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayName string, chatroomPath string) error {

	createTableIfNotExistsCRDB(ctx, client, "users", "(chatroompath STRING, displayname STRING)")

	_, err := client.Exec(ctx, "INSERT INTO users (chatroompath, displayname) VALUES ($1, $2)", chatroomPath, displayName)
	if err != nil {
		log.Println("Error adding user to chatroom:", err)
	}
	return err

}

func getAllUsersInChatroomCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) ([]string, error) {

	var users []string

	rows, err := client.Query(ctx, "SELECT displayname FROM users WHERE chatroompath = $1", chatroomPath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}

func isUserInChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayname string, chatroomPath string) (bool, error) {

	createTableIfNotExistsCRDB(ctx, client, "users", "(chatroompath STRING, displayname STRING)")

	var isMember bool
	err := client.QueryRow(ctx, "SELECT TRUE FROM users WHERE displayname = $1 AND chatroompath = $2", displayname, chatroomPath).Scan(&isMember)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	log.Println(isMember)
	return isMember, nil

}

func removeUserFromChatroomCRDB(ctx context.Context, client *pgxpool.Pool, displayName string, chatroomPath string) error {

	_, err := client.Exec(ctx, "DELETE FROM users WHERE displayname = $1 AND chatroompath = $2", displayName, chatroomPath)
	if err != nil {
		log.Println("Error removing user from chatroom:", err)
	}
	return err

}

func changeUserNameCRDB(ctx context.Context, client *pgxpool.Pool, oldName string, newName string, chatroomPath string) error {

	_, err := client.Exec(ctx, "UPDATE users SET displayname = $1 WHERE displayname = $2 AND chatroompath = $3", newName, oldName, chatroomPath)
	if err != nil {
		log.Println("Error updating user displayname:", err)
	}
	return err

}

func getMessageHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) ([]string, error) {

	createTableIfNotExistsCRDB(ctx, client, "messages", "(chatroompath STRING, message STRING[], CONSTRAINT chatroompath_unique UNIQUE (chatroompath))")

	var messages []string
	err := client.QueryRow(ctx, "SELECT message FROM messages WHERE chatroompath = $1", chatroomPath).Scan(&messages)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []string{}, nil
		} else {
			return nil, err
		}
	}
	log.Println(messages)
	return messages, nil

}

func addMessageToHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string, chatMessage string) error {

	queryString := "INSERT INTO messages (chatroompath, message) VALUES ($1, ARRAY[$2]) ON CONFLICT (chatroompath) DO update SET message = messages.message || EXCLUDED.message"

	_, err := client.Exec(ctx, queryString, chatroomPath, chatMessage)
	if err != nil {
		log.Println("Error adding message to chatroom history:", err)
	}
	return err

}

func removeMessageFromHistoryCRDB(ctx context.Context, client *pgxpool.Pool, chatroomPath string) error {

	_, err := client.Exec(
		ctx, `
		UPDATE messages
		SET message = (
			SELECT ARRAY_AGG(msg)
			FROM (
				SELECT UNNEST(message) AS msg
				FROM messages
				WHERE chatroompath = $1
				OFFSET 1
			)
		)
		WHERE chatroompath = $1;
	`,
		chatroomPath)
	if err != nil {
		log.Println("Error removing message from chatroom history:", err)
	}
	return err

}

func deleteKeyCRDB(ctx context.Context, client *pgxpool.Pool, key string) error {

	_, err := client.Exec(ctx, "DELETE FROM messages WHERE chatroompath = $1", key)
	if err != nil {
		log.Println("Error deleting message history of chatroom:", err)
	}
	return err

}

func getAllChatroomsCRDB(ctx context.Context, client *pgxpool.Pool, key string) (map[string]string, error) {

	createTableIfNotExistsCRDB(ctx, client, "lobby", "(roomname STRING, roompath STRING)")

	rooms := make(map[string]string)

	rows, err := client.Query(ctx, "SELECT * FROM lobby")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		rooms[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil

}

func addChatroomToLobbyCRDB(ctx context.Context, client *pgxpool.Pool, roomname string, roompath string) error {

	_, err := client.Exec(ctx, "INSERT INTO lobby (roomname, roompath) VALUES ($1, $2)", roomname, roompath)
	if err != nil {
		log.Println("Error adding chatroom to lobby:", err)
	}
	return err

}

func removeChatroomFromLobbyCRDB(ctx context.Context, client *pgxpool.Pool, roomname string) error {

	_, err := client.Exec(ctx, "DELETE FROM lobby WHERE roomname = $1", roomname)
	if err != nil {
		log.Println("Error removing chatroom from lobby:", err)
	}
	return err

}
