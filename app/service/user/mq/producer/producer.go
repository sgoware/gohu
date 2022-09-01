package producer

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type ChangeFollowerMessage struct {
	UserId int64 `json:"userId"`
	Action int32 `json:"action"`
}

func ChangeFollower(producer *nsq.Producer,
	rawMessage ChangeFollowerMessage) (err error) {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("marshal message filaed, %v", err)
	}

	err = producer.Publish("user-subject", message)
	if err != nil {
		return fmt.Errorf("publish msg to nsq failed, %v", err)
	}
	return nil
}
