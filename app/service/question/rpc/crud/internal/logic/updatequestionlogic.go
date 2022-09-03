package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/question/dao/model"
	modelpb "main/app/service/question/dao/pb"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"net/http"
	"time"
)

type UpdateQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionLogic {
	return &UpdateQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateQuestionLogic) UpdateQuestion(in *pb.UpdateQuestionReq) (res *pb.UpdateQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	nowTime := time.Now()

	questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject
	questionContentModel := l.svcCtx.QuestionModel.QuestionContent

	questionSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("question_subject_%d", in.QuestionId)).Bytes()
	if err == nil {
		questionSubjectProto := &modelpb.QuestionSubject{}
		err = proto.Unmarshal(questionSubjectBytes, questionSubjectProto)
		if err != nil {
			logger.Errorf("unmarshal [questionSubjectProto] failed, err: %v", err)
		} else {
			questionSubjectProto.Tag = in.Tag
			questionSubjectProto.Title = in.Title
			questionSubjectProto.Topic = in.Topic
			questionSubjectProto.UpdateTime = nowTime.String()

			questionSubjectBytes, err = proto.Marshal(questionSubjectProto)
			if err != nil {
				logger.Errorf("marshal [questionSubjectProto] failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			err = l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("question_subject_%d", in.QuestionId),
				questionSubjectBytes,
				time.Second*86400).Err()
			if err != nil {
				logger.Errorf("set [question_subject] cache failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}

			payload, err := json.Marshal(job.MsgCrudQuestionSubjectRecordPayload{
				Action:     2,
				Id:         in.QuestionId,
				Title:      in.Title,
				Topic:      in.Topic,
				Tag:        in.Tag,
				UpdateTime: nowTime,
			})
			if err != nil {
				logger.Debugf("marshal [MsgCrudQuestionSubjectRecordPayload] failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}

			_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgCrudQuestionSubjectRecordTask, payload))
			if err != nil {
				logger.Debugf("create [sgCrudQuestionSubjectRecordTask] insert queue failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
		}
	} else {
		if err != redis.Nil {
			logger.Errorf("get [question_subject] cache failed, err: %v", err)
		}
		_, err = questionSubjectModel.WithContext(l.ctx).Select(
			questionSubjectModel.ID,
			questionSubjectModel.Title,
			questionSubjectModel.Topic,
			questionSubjectModel.Tag).
			Where(questionSubjectModel.ID.Eq(in.QuestionId)).
			Updates(model.QuestionSubject{
				Title: in.Title,
				Topic: in.Topic,
				Tag:   in.Tag,
			})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				res = &pb.UpdateQuestionRes{
					Code: http.StatusForbidden,
					Msg:  "question not found",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			} else {
				logger.Errorf("update [question_subject] record failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
		}

		questionSubject, err := questionSubjectModel.WithContext(l.ctx).
			Where(questionSubjectModel.ID.Eq(in.QuestionId)).
			First()
		if err != nil {
			logger.Errorf("query [question_subject] record failed, err: %v", err)
		} else {
			questionSubjectBytes, err := proto.Marshal(&modelpb.QuestionSubject{
				Id:          questionSubject.ID,
				UserId:      questionSubject.UserID,
				IpLoc:       questionSubject.IPLoc,
				Title:       questionSubject.Title,
				Topic:       questionSubject.Topic,
				Tag:         questionSubject.Tag,
				SubCount:    questionSubject.SubCount,
				AnswerCount: questionSubject.AnswerCount,
				ViewCount:   questionSubject.ViewCount,
				State:       questionSubject.State,
				CreateTime:  questionSubject.CreateTime.String(),
				UpdateTime:  questionSubject.UpdateTime.String(),
			})
			if err != nil {
				logger.Errorf("marshal [questionSubjectProto] failed, err: %v", err)
			} else {
				err = l.svcCtx.Rdb.Set(l.ctx,
					fmt.Sprintf("question_subject_%d", in.QuestionId),
					questionSubjectBytes,
					time.Second*86400).Err()
				if err != nil {
					logger.Errorf("set [question_subject] cache failed, err: %v", err)
				}
			}
		}
	}

	questionContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("question_content_%d", in.QuestionId)).Bytes()
	if err == nil {
		questionContentProto := &modelpb.QuestionContent{}
		err = proto.Unmarshal(questionContentBytes, questionContentProto)
		if err != nil {
			logger.Errorf("unmarshal [questionSubjectProto] failed, err: %v", err)
		} else {
			questionContentProto.Content = in.Content
			questionContentProto.UpdateTime = nowTime.String()

			questionContentBytes, err = proto.Marshal(questionContentProto)
			if err != nil {
				logger.Errorf("marshal [questionContentProto] failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal er",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			err = l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("question_content_%d", in.QuestionId),
				questionContentBytes,
				time.Second*86400).Err()
			if err != nil {
				logger.Errorf("set [question_content] cache failed, err: %v", err)
				res = &pb.UpdateQuestionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal er",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			res = &pb.UpdateQuestionRes{
				Code: http.StatusOK,
				Msg:  "update question successfully",
				Ok:   true,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	} else {
		if err != redis.Nil {
			logger.Errorf("get [question_content] cache failed, err: %v", err)
		}
		_, err = questionContentModel.WithContext(l.ctx).Select(
			questionContentModel.QuestionID,
			questionContentModel.Content).
			Where(questionContentModel.QuestionID.Eq(in.QuestionId)).
			Update(questionContentModel.Content, in.Content)
		if err != nil {
			logger.Errorf("update [question_content] record failed, err: %v", err)
			res = &pb.UpdateQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}

		questionContent, err := questionContentModel.WithContext(l.ctx).
			Where(questionContentModel.QuestionID.Eq(in.QuestionId)).
			First()
		if err != nil {
			logger.Errorf("query [question_content] record failed, err: %v", err)
		} else {
			questionContentBytes, err := proto.Marshal(&modelpb.QuestionContent{
				QuestionId: questionContent.QuestionID,
				Content:    questionContent.Content,
				Meta:       questionContent.Meta,
				CreateTime: questionContent.CreateTime.String(),
				UpdateTime: questionContent.UpdateTime.String(),
			})
			if err != nil {
				logger.Errorf("marshal [questionContentProto] failed, err: %v", err)
			} else {
				err = l.svcCtx.Rdb.Set(l.ctx,
					fmt.Sprintf("question_content_%d", in.QuestionId),
					questionContentBytes,
					time.Second*86400).Err()
				if err != nil {
					logger.Errorf("set [question_content] cache failed, err: %v", err)
				}
			}
		}

	}

	res = &pb.UpdateQuestionRes{
		Code: http.StatusOK,
		Msg:  "update question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
