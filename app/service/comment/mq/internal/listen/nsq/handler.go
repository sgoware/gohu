package nsq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"main/app/service/comment/rpc/crud/crud"
	questionMqProduct "main/app/service/question/rpc/crud/mq"
)

type CommentSubjectHandler struct {
	CrudRpcClient crud.Crud
}

func (m *CommentSubjectHandler) HandleMessage(nsqMsg *nsq.Message) (err error) {
	msg := &questionMqProduct.InitMessage{}
	err = json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return err
	}
	switch msg.Action {
	case "init":
		res, _ := m.CrudRpcClient.InitSubject(context.Background(), &crud.InitSubjectReq{
			ObjType: msg.Data.ObjType,
			ObjId:   msg.Data.ObjId,
		})
		if !res.Ok {
			return errors.New(res.Msg)
		}

	case "delete":
		res, _ := m.CrudRpcClient.DeleteSubject(context.Background(), &crud.DeleteSubjectReq{
			ObjType: msg.Data.ObjType,
			ObjId:   msg.Data.ObjId,
		})
		if !res.Ok {
			return errors.New(res.Msg)
		}
	}
	return nil
}
