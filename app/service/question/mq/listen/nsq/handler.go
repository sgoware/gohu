package nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	questionMqProducer "main/app/service/question/mq/producer"
	"main/app/service/question/rpc/crud/crud"
)

type ChangeAttrHandler struct {
	QuestionRpcClient crud.Crud
}

func (m *ChangeAttrHandler) HandleMessage(nsqMsg *nsq.Message) (err error) {
	msg := &questionMqProducer.CollectMessage{}
	err = json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("unmarshal msg failed, %v", err)
	}

	ctx := context.Background()
	res, _ := m.QuestionRpcClient.ChangeAttr(ctx, &crud.ChangeAttrReq{
		ObjType:  msg.ObjType,
		ObjId:    msg.ObjId,
		AttrType: msg.AttrType,
		Action:   msg.Action,
	})
	if !res.Ok {
		return fmt.Errorf("change attr failed, %v", res.Msg)
	}

	return nil
}
