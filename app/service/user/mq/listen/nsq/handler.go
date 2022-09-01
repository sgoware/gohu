package nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	userMqProducer "main/app/service/user/mq/producer"
	"main/app/service/user/rpc/crud/crud"
)

type ChangeFollowerHandler struct {
	CrudRpcClient crud.Crud
}

func (m *ChangeFollowerHandler) HandleMessage(nsqMsg *nsq.Message) (err error) {
	msg := &userMqProducer.ChangeFollowerMessage{}
	err = json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("unmarshal msg failed, %v", err)
	}
	res, _ := m.CrudRpcClient.ChangeFollower(context.Background(), &crud.ChangeFollowerReq{
		UserId: msg.UserId,
		Action: msg.Action,
	})
	if !res.Ok {
		return fmt.Errorf("change follower failed, %v", res.Msg)
	}
	return nil
}
