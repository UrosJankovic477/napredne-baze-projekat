package internals

import (
	"context"
	"fmt"
	"time"

	"github.com/dranikpg/gtrs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/philippseith/signalr"
	"github.com/redis/go-redis/v9"
)

var Server signalr.Server

func GetHubClients() signalr.HubClients {
	return Server.HubClients()
}

type ChatHub struct {
	signalr.Hub
}

func (h *ChatHub) OnConnected(connectionID string) {
	fmt.Printf("%s connected\n", connectionID)
}

func (h *ChatHub) OnDisconnected(connectionID string) {
	fmt.Printf("%s disconnected\n", connectionID)
}

func (h *ChatHub) SendMessage(token string, message string, UUID string) {
	//	fmt.Println(message)
	user_node, _, err := GetUserFromToken(token)
	if err != nil {
		return
	}
	username := user_node.Props["Username"]
	SendMessageToStream(Message{
		Username: username.(string),
		Content:  message,
		Time:     time.Now().Unix()},
		UUID)
}

func (h *ChatHub) JoinChat(UUID string) {
	h.Groups().AddToGroup(UUID, h.ConnectionID())
	msgs, err := GetMessages(UUID)
	if err != nil {
		return
	}
	h.Clients().Caller().Send("ReceiveMessageList", msgs)

}

var ctx context.Context
var driver neo4j.DriverWithContext
var rdb *redis.Client
var Hub ChatHub

func Initialize() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	dbUri := "neo4j://localhost"
	var err error
	driver, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbName, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	ctx = context.Background()
	_, err = doQuery("CREATE CONSTRAINT unique_username IF NOT EXISTS FOR (usr:AccountCredentials) REQUIRE usr.Username IS UNIQUE ", nil)
	if err != nil {
		panic(err)
	}
	_, err = doQuery("CREATE CONSTRAINT unique_forum_name IF NOT EXISTS FOR (forum:Forum) REQUIRE forum.Name IS UNIQUE ", nil)
	if err != nil {
		panic(err)
	}

	keys, err := rdb.HKeys(ctx, "chatrooms").Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		consumer := gtrs.NewConsumer[Message](ctx, rdb, gtrs.StreamIDs{key: "$"})
		go func(UUID string) {
			for {
				gtrs_msg := <-consumer.Chan()
				msg := gtrs_msg.Data
				GetHubClients().Group(UUID).Send("ReceiveMessage", msg.Username, msg.Content, msg.Time)
			}
		}(key)
	}

	Server, _ = signalr.NewServer(context.Background(),
		signalr.SimpleHubFactory(&Hub))

}
