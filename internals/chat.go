package internals

import (
	"errors"
	"sync"

	"github.com/dranikpg/gtrs"
)

type Message struct {
	Username string
	Content  string
	Time     int64
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
