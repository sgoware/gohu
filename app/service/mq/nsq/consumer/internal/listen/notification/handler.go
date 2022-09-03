package notification

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
	"sync"
)

type PublishNotificationHandler struct {
	Domain                    string
	NotificationCrudRpcClient crud.Crud
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
		data := &notificationMqProducer.FollowerData{}

		bytesData, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("marshal msg data failed, %v", err)
		}

		err = json.Unmarshal(bytesData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal msg data failed, %v", err)
		}

		userInfoRes, err := req.NewRequest().Get(
			fmt.Sprintf("https://%s/api/user/profile/%s", m.Domain, cast.ToString(data.FollowerId)))
		if err != nil {
			return fmt.Errorf("query user info failed, %v", err)
		}

		j := gjson.Parse(userInfoRes.String())
		if !j.Get("ok").Bool() {
			return fmt.Errorf("query user info failed, %v", j.Get("msg").String())
		}

		rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
			UserId:      data.UserId,
			MessageType: 1,
			Title:       fmt.Sprintf("用户 %s 关注了你", j.Get("data.nickname").String()),
			Content:     "",                                                              // 空
			Url:         fmt.Sprintf("https://%s/profile/%d", m.Domain, data.FollowerId), // 用户主页
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
			answerRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/question/answer/%d", m.Domain, data.ObjId))
			if err != nil {
				return fmt.Errorf("query answer info failed, %v", err)
			}
			answerJson := gjson.Parse(answerRes.String())

			questionRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/question/question/%d",
					m.Domain, answerJson.Get("data.answer_index.question_id").Int()))
			if err != nil {
				return fmt.Errorf("query question info failed, %v", err)
			}
			questionJson := gjson.Parse(questionRes.String())

			userInfoRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/user/profile/%d", m.Domain, data.UserId))
			if err != nil {
				return fmt.Errorf("query user info failed, %v", err)
			}

			userInfoJson := gjson.Parse(userInfoRes.String())
			if !userInfoJson.Get("ok").Bool() {
				return fmt.Errorf("query user info failed, %v", userInfoJson.Get("msg"))
			}

			var title string
			if data.Action == 1 {
				title = fmt.Sprintf("用户 %s 赞同了在问题 %s 下的回答",
					userInfoJson.Get("data.nickname").String(),
					questionJson.Get("data.question_subject.title"))
			} else {
				title = fmt.Sprintf("用户 %s 喜欢了在问题 %s 下的回答",
					userInfoJson.Get("data.nickname").String(),
					questionJson.Get("data.question_subject.title"))
			}
			rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
				UserId:      answerJson.Get("data.answer_index.user_id").Int(),
				MessageType: 2,
				Title:       title,
				Content:     "",
				Url: fmt.Sprintf("https://%s/question/%s/answer/%s",
					m.Domain,
					questionJson.Get("data.question_subject.id").String(),
					answerJson.Get("data.answer_index.id").String()),
			})
			if !rpcRes.Ok {
				return fmt.Errorf("publish notification failed, %v", rpcRes.Msg)
			}

		case 2:
			// 文章

		case 3:
			// 评论
			userInfoRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/user/profile/%d", m.Domain, data.UserId))
			if err != nil {
				return fmt.Errorf("query user info failed, err: %v", err)
			}

			userInfoJson := gjson.Parse(userInfoRes.String())
			if !userInfoJson.Get("ok").Bool() {
				return fmt.Errorf("query user info failed, err: %v",
					userInfoJson.Get("msg"))
			}

			nickname := userInfoJson.Get("data.nickname")

			commentInfoRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/comment/%d", m.Domain, data.ObjId))
			if err != nil {
				return fmt.Errorf("query comment info failed, err: %v", err)
			}

			commentInfoJson := gjson.Parse(commentInfoRes.String())
			if !commentInfoJson.Get("ok").Bool() {
				return fmt.Errorf("query comment info failed, err: %v",
					commentInfoJson.Get("msg"))
			}

			userId := commentInfoJson.Get("data.comment_index.user_id").Int()

			rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
				UserId:      userId,
				MessageType: 2,
				Title:       fmt.Sprintf("用户 %s 赞同了你的评论", nickname),
				Content:     "",
				Url:         "",
			})
			if !rpcRes.Ok {
				return fmt.Errorf("publish notification failed, %v", rpcRes.Msg)
			}

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
				fmt.Sprintf("https://%s/api/comment/subject/%d", m.Domain, data.SubjectId))
			if err != nil {
				return fmt.Errorf("query comment subject failed, err: %v", err)
			}
			commentSubjectJson := gjson.Parse(commentSubjectRes.String())
			objType := commentSubjectJson.Get("data.comment_subject.obj_type").Int()
			objId := commentSubjectJson.Get("data.comment_subject.obj_id").Int()
			switch objType {
			case 1:
				// 评论回答
				answerRes, err := req.NewRequest().Get(
					fmt.Sprintf("https://%s/api/question/answer/%d", m.Domain, objId))
				if err != nil {
					return fmt.Errorf("query answer failed, err: %v", err)
				}
				answerJson := gjson.Parse(answerRes.String())
				userId := answerJson.Get("data.answer_index.user_id").Int()
				questionId := answerJson.Get("data.answer_index.question_id").Int()

				questionRes, err := req.NewRequest().Get(
					fmt.Sprintf("https://%s/api/question/question/%d", m.Domain, questionId))
				if err != nil {
					return fmt.Errorf("query question failed, err: %v", err)
				}
				questionJson := gjson.Parse(questionRes.String())
				questionTitle := questionJson.Get("data.question_subject.title").String()

				userInfoRes, err := req.NewRequest().Get(
					fmt.Sprintf("https://%s/api/user/profile/%d", m.Domain, data.UserId))
				if err != nil {
					return fmt.Errorf("query user info failed, err: %v", err)
				}
				userInfoJson := gjson.Parse(userInfoRes.String())
				userNickName := userInfoJson.Get("data.nickname").String()

				rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
					UserId:      userId,
					MessageType: 3,
					Title:       fmt.Sprintf("用户 %s 评论了你的回答 %s", userNickName, questionTitle),
					Content:     "",
					Url:         fmt.Sprintf("https://%s/question/%d/answer/%d", m.Domain, questionId, objId),
				})
				if !rpcRes.Ok {
					return fmt.Errorf("publish comment notification failed, %v", rpcRes.Msg)
				}

			case 2:
			// 文章

			default:

			}

		} else {
			// 回复评论
			commentInfoRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/comment/%d", m.Domain, data.CommentId))
			if err != nil {
				return fmt.Errorf("query [commentIndex] failed, err: %v", err)
			}

			commentInfoJson := gjson.Parse(commentInfoRes.String())

			if !commentInfoJson.Get("ok").Bool() {
				return fmt.Errorf("query [commentIndex] failed, err: %v",
					commentInfoJson.Get("msg"))
			}

			userInfoRes, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/user/profile/%d",
					m.Domain,
					data.UserId))
			if err != nil {
				return fmt.Errorf("query [userInfo] failed, err: %v", err)
			}

			userInfoResJson := gjson.Parse(userInfoRes.String())
			if !userInfoResJson.Get("ok").Bool() {
				return fmt.Errorf("query [userInfo] failed, err: %v", userInfoResJson.Get("msg"))
			}

			rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
				UserId:      commentInfoJson.Get("data.comment_index.user_id").Int(),
				MessageType: 3,
				Title: fmt.Sprintf("用户 %s 回复了你的评论",
					userInfoResJson.Get("data.nickname"),
				),
				Content: "",
				Url:     "",
			})

			if !rpcRes.Ok {
				return fmt.Errorf("publish comment notification failed, %v", rpcRes.Msg)
			}
		}

	case 4:
		// 关注的人
		data := &notificationMqProducer.SubscriptionData{}

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
			// 发布问题
			questionInfoResJson, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/question/question/%d", m.Domain, data.ObjId))
			if err != nil {
				return fmt.Errorf("query question info failed, err: %v", err)
			}

			questionInfoRes := gjson.Parse(questionInfoResJson.String())
			if !questionInfoRes.Get("ok").Bool() {
				return fmt.Errorf("query question info failed, err: %v",
					questionInfoRes.Get("msg").String())
			}

			userInfoResJson, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/user/profile/%s", m.Domain, cast.ToString(data.UserId)))
			if err != nil {
				return fmt.Errorf("query user info failed, err: %v", err)
			}

			userInfoRes := gjson.Parse(userInfoResJson.String())
			if !userInfoRes.Get("ok").Bool() {
				return fmt.Errorf("query user info failed, err: %v",
					userInfoRes.Get("msg").String())
			}

			followerResJson, err := req.NewRequest().Get(
				fmt.Sprintf("https://%s/api/user/follower/%d",
					m.Domain,
					data.UserId))
			if err != nil {
				return fmt.Errorf("query followers failed, err: %v", err)
			}

			followersRes := gjson.Parse(followerResJson.String())
			if !followersRes.Get("ok").Bool() {
				return fmt.Errorf("query follower failed, err: %v",
					followersRes.Get("msg").String())
			}

			followers := followersRes.Get("data.user_ids").Array()

			wg := sync.WaitGroup{}

			wg.Add(len(followers))

			for _, follower := range followers {
				go func(follower gjson.Result) {
					defer wg.Done()

					_, _ = m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
						UserId:      follower.Int(),
						MessageType: 4,
						Title: fmt.Sprintf("你关注的用户 %s 提出了问题 %s",
							userInfoRes.Get("data.nickname"),
							questionInfoRes.Get("data.question_subject.title")),
						Content: "",
						Url:     fmt.Sprintf("https://%s/question/%d", m.Domain, data.ObjId),
					})
				}(follower)
			}

			wg.Wait()

			return nil
		}

	case 5:
		// 问题回答
		data := &notificationMqProducer.AnswerData{}

		bytesData, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("marshal msg data failed, %v", err)
		}

		err = json.Unmarshal(bytesData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal msg data failed, %v", err)
		}

		questionInfoResJson, err := req.NewRequest().Get(
			fmt.Sprintf("https://%s/api/question/question/%d", m.Domain, data.QuestionId))
		if err != nil {
			return fmt.Errorf("query question info failed, err: %v", err)
		}

		questionInfoRes := gjson.Parse(questionInfoResJson.String())
		if !questionInfoRes.Get("ok").Bool() {
			return fmt.Errorf("query question info failed, err: %v", questionInfoRes.Get("msg"))
		}

		userInfoResJson, err := req.NewRequest().Get(
			fmt.Sprintf("https://%s/api/user/profile/%s", m.Domain, cast.ToString(data.UserId)))
		if err != nil {
			return fmt.Errorf("query user info failed, err: %v", err)
		}

		userInfoRes := gjson.Parse(userInfoResJson.String())
		if !userInfoRes.Get("ok").Bool() {
			return fmt.Errorf("query user info failed, err: %v", userInfoRes.Get("msg"))
		}

		rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
			UserId:      questionInfoRes.Get("data.question_subject.user_id").Int(),
			MessageType: 5,
			Title: fmt.Sprintf("用户 %s 回答了你的问题 %s",
				userInfoRes.Get("data.nickname"),
				questionInfoRes.Get("data.question_subject.title")),
			Content: "",
			Url: fmt.Sprintf("https://%s/question/%d/answer/%d",
				m.Domain,
				data.QuestionId,
				data.AnswerId),
		})
		if !rpcRes.Ok {
			return fmt.Errorf("publish subscription notification failed, err: %v", rpcRes.Msg)
		}

	}
	return nil
}
