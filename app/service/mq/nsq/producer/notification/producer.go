package notification

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type PublishNotificationMessage struct {
	MessageType int32       `json:"message_type"`
	Data        interface{} `json:"data"`
}

type FollowerData struct {
	UserId     int64 `json:"user_id"` // 被通知的用户
	FollowerId int64 `json:"follower_id"`
}

type ApproveAndLikeData struct {
	UserId  int64 `json:"user_id"`
	Action  int32 `json:"action"`
	ObjType int32 `json:"obj_type"`
	ObjId   int64 `json:"obj_id"`
}

type CommentData struct {
	UserId    int64 `json:"user_id"`
	SubjectId int64 `json:"subject_id"`
	CommentId int64 `json:"comment_id"`
}

type SubscriptionData struct {
	UserId  int64 `json:"user_id"`  // 被关注者id
	Action  int32 `json:"action"`   // 1: 创建 2: 更新
	ObjType int32 `json:"obj_type"` // 1: 发布问题
	ObjId   int64 `json:"obj_id"`
}

type AnswerData struct {
	UserId     int64 `json:"user_id"`     // 回答者用户id
	QuestionId int64 `json:"question_id"` // 回答的问题id
	AnswerId   int64 `json:"answer_id"`   // 回答id
}

func PublishNotification(producer *nsq.Producer,
	rawMessage PublishNotificationMessage) (err error) {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("marshal message filaed, %v", err)
	}

	err = producer.Publish("notification-publish", message)
	if err != nil {
		return fmt.Errorf("publish msg to nsq failed, %v", err)
	}
	return nil
}
