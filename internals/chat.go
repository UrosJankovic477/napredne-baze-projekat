package internals

import (
	"errors"
	"net/http"
	"sync"

	"github.com/dranikpg/gtrs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Message struct {
	Username string
	Content  string
	Time     int64
}

type Chatroom struct {
	UUID  string
	Name  string
	Users []string
}

var Mutex sync.Mutex

func MakeChatRoom() string {
	Mutex.Lock()
	UUID := GenerateUUID()
	stream := gtrs.NewStream[Message](rdb, UUID, &gtrs.Options{TTL: gtrs.NoExpiration, MaxLen: 10000, Approx: true})

	rdb.HSet(ctx, "chatrooms", UUID, true).Result()

	consumer := gtrs.NewConsumer[Message](ctx, rdb, gtrs.StreamIDs{UUID: "$"})
	go func() {
		for {
			gtrs_msg := <-consumer.Chan()
			msg := gtrs_msg.Data
			GetHubClients().Group(UUID).Send("ReceiveMessage", msg.Username, msg.Content, msg.Time)
		}
	}()

	Mutex.Unlock()
	return stream.Key()
}

func GetMessages(UUID string) ([]Message, error) {
	Mutex.Lock()

	ok, err := rdb.HExists(ctx, "chatrooms", UUID).Result()
	if !ok || err != nil {
		Mutex.Unlock()
		return nil, errors.New("Chat room not found")
	}
	stream := gtrs.NewStream[Message](rdb, UUID, &gtrs.Options{TTL: gtrs.NoExpiration, MaxLen: 10000, Approx: true})
	gtrs_msgs, err := stream.Range(ctx, "-", "+")

	msgs := make([]Message, 0)

	for _, msg := range gtrs_msgs {
		msgs = append(msgs, msg.Data)
	}

	Mutex.Unlock()

	return msgs, err
}

func SendMessageToStream(msg Message, UUID string) error {
	Mutex.Lock()

	ok, err := rdb.HExists(ctx, "chatrooms", UUID).Result()
	if !ok || err != nil {
		Mutex.Unlock()
		return errors.New("Chat room not found")
	}
	stream := gtrs.NewStream[Message](rdb, UUID, &gtrs.Options{TTL: gtrs.NoExpiration, MaxLen: 10000, Approx: true})

	stream.Add(ctx, msg)

	Mutex.Unlock()
	return nil
}

func GetUsersChatrooms(token string) ([]Chatroom, int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}
	chatrooms_node, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (chatroom:Chatroom) WHERE (usr)-[:IN_CHATROOM]-(chatroom) "+
		"WITH chatroom MATCH (usr) WHERE (usr)-[:IN_CHATROOM]-(chatroom) "+
		"RETURN chatroom, usr",
		map[string]any{
			"Id": user_node.ElementId,
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	chatroom_map := make(map[string]Chatroom, 0)

	for _, record := range chatrooms_node.Records {
		chatroom_node, _ := record.Get("chatroom")
		usr_node, _ := record.Get("usr")
		uuid := chatroom_node.(dbtype.Node).Props["UUID"].(string)
		name := chatroom_node.(dbtype.Node).Props["Name"].(string)
		username := usr_node.(dbtype.Node).Props["Username"].(string)
		chatroom, ok := chatroom_map[uuid]
		if !ok {
			chatroom = Chatroom{UUID: uuid, Name: name, Users: make([]string, 0)}
		}
		chatroom.Users = append(chatroom.Users, username)
		chatroom_map[uuid] = chatroom
	}
	chatrooms := make([]Chatroom, 0)
	for _, chatroom := range chatroom_map {
		chatrooms = append(chatrooms, chatroom)
	}

	return chatrooms, http.StatusOK, nil
}
