package producer

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type CollectMessage struct {
	ObjType  int32 `json:"objType"`
	ObjId    int64 `json:"objId"`
	AttrType int32 `json:"attrType"`
	Action   int32 `json:"action"`
}

func DoCollect(producer *nsq.Producer,
	rawMessage CollectMessage) (err error) {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("marshal message filaed, %v", err)
	}

	err = producer.Publish("question-collect", message)
	if err != nil {
		return fmt.Errorf("publish msg to nsq failed, %v", err)
	}
	return nil
}
