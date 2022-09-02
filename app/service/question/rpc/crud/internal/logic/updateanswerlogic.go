package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/app/common/log"
	modelpb "main/app/service/question/dao/pb"
	"net/http"
	"time"

	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnswerLogic {
	return &UpdateAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAnswerLogic) UpdateAnswer(in *pb.UpdateAnswerReq) (res *pb.UpdateAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	nowTime := time.Now()

	answerContentModel := l.svcCtx.QuestionModel.AnswerContent

	answerContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("answer_content", in.AnswerId)).Bytes()
	if err == nil {
		answerContentProto := &modelpb.AnswerContent{}
		err = proto.Unmarshal(answerContentBytes, answerContentProto)
		if err != nil {
			logger.Errorf("unmarshal [answerContentProto] failed, err: %v", err)
		} else {
			answerContentProto.AnswerId = in.AnswerId
			answerContentProto.Content = in.Content
			answerContentProto.UpdateTime = nowTime.String()

			answerContentBytes, err = proto.Marshal(answerContentProto)
			if err != nil {
				logger.Errorf("marshal [answerContentProto] failed, err: %v", err)
				res = &pb.UpdateAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			err = l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("answer_content_%d", in.AnswerId),
				answerContentBytes,
				time.Second*86400).Err()
			if err != nil {
				logger.Errorf("set [answer_content] cache failed, err: %v", err)
				res = &pb.UpdateAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %d", err)
				return res, nil
			}
		}
	} else {
		if err != redis.Nil {
			logger.Errorf("get [answer_content] cache failed, err: %v", err)
		}
		_, err = answerContentModel.WithContext(l.ctx).
			Where(answerContentModel.AnswerID.Eq(in.AnswerId)).
			Update(answerContentModel.Content, in.Content)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				res = &pb.UpdateAnswerRes{
					Code: http.StatusForbidden,
					Msg:  "answer not found",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			} else {
				logger.Errorf("update [answer_content] record failed, err: %v", err)
				res = &pb.UpdateAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
		}

		answerContent, err := answerContentModel.WithContext(l.ctx).
			Where(answerContentModel.AnswerID.Eq(in.AnswerId)).
			First()
		if err != nil {
			logger.Errorf("query [ansewr_content] record failed, err: %v", err)
		} else {
			answerContentBytes, err := proto.Marshal(&modelpb.AnswerContent{
				AnswerId:   answerContent.AnswerID,
				Content:    answerContent.Content,
				Meta:       answerContent.Meta,
				CreateTime: answerContent.CreateTime.String(),
				UpdateTime: answerContent.UpdateTime.String(),
			})
			if err != nil {
				logger.Errorf("marshal [answerContentProto] failed, err: %v", err)
			} else {
				err = l.svcCtx.Rdb.Set(l.ctx,
					fmt.Sprintf("answer_content_%d", in.AnswerId),
					answerContentBytes,
					time.Second*86400).Err()
				if err != nil {
					logger.Errorf("set [answer_content] cache failed, err: %v", err)
				}
			}
		}
	}

	res = &pb.UpdateAnswerRes{
		Code: http.StatusOK,
		Msg:  "update answer successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
