package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	notificationMqProducer "main/app/service/mq/nsq/producer/notification"
	"main/app/service/notification/rpc/crud/crud"
)

type PublishNotificationHandler struct {
	Domain        string
	CrudRpcClient crud.Crud
}

func (m *PublishNotificationHandler) HandleMessage(nsqMsg *nsq.Message) (err error) {
	msg := &notificationMqProducer.PublishNotificationMessage{}
	err = json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("unmarshal msg failed, %v", err)
	}

	ctx := context.Background()
	switch msg.MessageType {
	case 1:
		// 关注我的
		data := &notificationMqProducer.SubscriptionData{}

		bytesData, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("marshal msg data failed, %v", err)
		}

		err = json.Unmarshal(bytesData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal msg data failed, %v", err)
		}

		userInfoRes, err := req.NewRequest().Get(m.Domain + "/api/user/profile/" + cast.ToString(data.FollowerId))
		if err != nil {
			return fmt.Errorf("query user info failed, %v", err)
		}
		j := gjson.Parse(userInfoRes.String())
		if j.Get("ok").Bool() {
			return fmt.Errorf("query user info failed, %v", j.Get("msg").String())
		}

		rpcRes, _ := m.CrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
			UserId:      data.UserId,
			MessageType: 1,
			Title:       fmt.Sprintf("用户 %s 关注了你", j.Get("data.nickname").String()),
			Content:     "",                                                      // 空
			Url:         fmt.Sprintf("%s/profile/%d", m.Domain, data.FollowerId), // 用户主页
		})
		if !rpcRes.Ok {
			return fmt.Errorf("publish notification failed, %v", rpcRes.Msg)
		}

	case 2:
		// 赞同与喜欢
		data := &notificationMqProducer.ApproveAndLikeData{}

		bytesData, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("marshal msg data failed, %v", err)
		}

		err = json.Unmarshal(bytesData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal msg data failed, %v", err)
		}
		switch data.ObjType {
		case 1:
			// 回答
			answerRes, err := req.NewRequest().Get(fmt.Sprintf("%s/api/answer/%d",
				m.Domain, data.ObjId))
			if err != nil {
				return fmt.Errorf("query answer info failed, %v", err)
			}
			answerJson := gjson.Parse(answerRes.String())

			questionRes, err := req.NewRequest().Get(fmt.Sprintf("%s/api/question/%d",
				m.Domain, answerJson.Get("data.answer_index.question_id").Int()))
			if err != nil {
				return fmt.Errorf("query question info failed, %v", err)
			}
			questionJson := gjson.Parse(questionRes.String())

			userInfoRes, err := req.NewRequest().Get(m.Domain + "/api/user/profile/" + cast.ToString(data.UserId))
			if err != nil {
				return fmt.Errorf("query user info failed, %v", err)
			}
			userInfoJson := gjson.Parse(userInfoRes.String())
			if userInfoJson.Get("ok").Bool() {
				return fmt.Errorf("query user info failed, %v", answerJson.Get("msg").String())
			}

			var title string
			if data.Action == 1 {
				title = fmt.Sprintf("用户 %s 喜欢了在问题 %s 下的回答",
					userInfoJson.Get("data.nickname").String(),
					answerJson.Get("data.question_subject.title"))
			} else {
				title = fmt.Sprintf("用户 %s 赞同了在问题 %s 下的回答",
					userInfoJson.Get("data.nickname").String(),
					answerJson.Get("data.question_subject.title"))
			}
			rpcRes, _ := m.CrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
				UserId:      answerJson.Get("data.answer_index.id").Int(),
				MessageType: 2,
				Title:       title,
				Content:     "",
				Url: fmt.Sprintf("%s/question/%s/answer/%s",
					m.Domain,
					questionJson.Get("data.question_subject.id").String(),
					answerJson.Get("data.answer_index.id").String()),
			})
			if !rpcRes.Ok {
				return fmt.Errorf("publish notification failed, %v", rpcRes.Msg)
			}

		case 2:
			// 文章
			// TODO: 通知系统: 文章赞同与喜欢
		}

	case 3:
		// 评论与回复
		data := &notificationMqProducer.CommentData{}

		bytesData, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("marshal msg data failed, %v", err)
		}

		err = json.Unmarshal(bytesData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal msg data failed, %v", err)
		}

		if data.CommentId == 0 {
			commentSubjectRes, err := req.NewRequest().Get(
				fmt.Sprintf("%s/api/comment/subject/%d", m.Domain, data.SubjectId))
			if err != nil {
				return fmt.Errorf("query comment subject failed, err: %v", err)
			}
			commentSubjectJson := gjson.Parse(commentSubjectRes.String())
			objType := commentSubjectJson.Get("data.obj_type").Int()
			objId := commentSubjectJson.Get("data.obj_id").Int()
			switch objType {
			case 1:
				// 回答
				answerRes, err := req.NewRequest().Get(
					fmt.Sprintf("%s/api/question/answer/%d", m.Domain, objId))
				if err != nil {
					return fmt.Errorf("query answer failed, err: %v", err)
				}
				answerJson := gjson.Parse(answerRes.String())
				userId := answerJson.Get("data.answer_index.user_id").Int()
				questionId := answerJson.Get("data.answer_index.question_id").Int()

				questionRes, err := req.NewRequest().Get(
					fmt.Sprintf("%s/api/question/question/%d", m.Domain, questionId))
				if err != nil {
					return fmt.Errorf("query question failed, err: %v", err)
				}
				questionJson := gjson.Parse(questionRes.String())
				questionTitle := questionJson.Get("data.question_subject.title").String()

				userInfoRes, err := req.NewRequest().Get(
					fmt.Sprintf("%s/api/user/profile/%d", m.Domain, data.UserId))
				if err != nil {
					return fmt.Errorf("query user info failed, err: %v", err)
				}
				userInfoJson := gjson.Parse(userInfoRes.String())
				userNickName := userInfoJson.Get("data.nickname").String()

				rpcRes, _ := m.CrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
					UserId:      userId,
					MessageType: 3,
					Title:       fmt.Sprintf("用户 %s 评论了你的回答 %s", userNickName, questionTitle),
					Content:     "",
					Url:         fmt.Sprintf("%s/question/%d/answer/%d", m.Domain, questionId, objId),
				})
				if !rpcRes.Ok {
					return fmt.Errorf("publish comment notification failed, %v", rpcRes.Msg)
				}

			case 2:
				// 文章
			}

		} else {

		}

	case 4:
		// 邀请
		// TODO: 通知系统: 邀请回答

	case 5:
		// 提到我的
		// TODO: 通知系统: 评论中提到我的
	}
	//msg := &questionMqProduce.AnswerSubjectMessage{}
	//err = json.Unmarshal(nsqMsg.Body, &msg)
	//if err != nil {
	//	return err
	//}
	//switch msg.Action {
	//case "init":
	//	res, _ := m.CrudRpcClient.InitSubject(context.Background(), &crud.InitSubjectReq{
	//		ObjType: msg.Data.ObjType,
	//		ObjId:   msg.Data.ObjId,
	//	})
	//	if !res.Ok {
	//		return errors.New(res.Msg)
	//	}
	//
	//case "delete":
	//	res, _ := m.CrudRpcClient.DeleteSubject(context.Background(), &crud.DeleteSubjectReq{
	//		ObjType: msg.Data.ObjType,
	//		ObjId:   msg.Data.ObjId,
	//	})
	//	if !res.Ok {
	//		return errors.New(res.Msg)
	//	}
	//}
	return nil
}
