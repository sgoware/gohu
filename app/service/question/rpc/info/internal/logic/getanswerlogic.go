package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/question/rpc/info/internal/svc"
	"main/app/service/question/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnswerLogic {
	return &GetAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAnswerLogic) GetAnswer(in *pb.GetAnswerReq) (res *pb.GetAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetAnswerRes_Data{
		AnswerIndex:   &pb.AnswerIndex{},
		AnswerContent: &pb.AnswerContent{},
	}

	answerIndexBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("answer_index_%d", in.AnswerId)).Bytes()
	if err == nil {
		answerIndexProto := &pb.AnswerIndex{}
		err = proto.Unmarshal(answerIndexBytes, answerIndexProto)
		if err != nil {
			logger.Errorf("unmarshal [answerIndexProto] failed, err: %v", err)
			res = &pb.GetAnswerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}

		approveCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("answer_index_approve_cnt_%d", in.AnswerId)).Int64()
		if err != nil {
			if err != redis.Nil {
				logger.Errorf("get [answer_index_approve_cnt] failed, err: %v", err)
				res = &pb.GetAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		} else {
			resData.AnswerIndex.ApproveCount += int32(approveCnt)
		}

		likeCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("answer_index_like_cnt_%d", in.AnswerId)).Int64()
		if err != nil {
			if err != redis.Nil {
				logger.Errorf("get [answer_index_like_cnt] failed, err: %v", err)
				res = &pb.GetAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		} else {
			resData.AnswerIndex.LikeCount += int32(likeCnt)
		}

		collectCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("answer_index_collect_cnt_%d", in.AnswerId)).Int64()
		if err != nil {
			if err != redis.Nil {
				logger.Errorf("get [answer_index_collect_cnt] failed, err: %v", err)
				res = &pb.GetAnswerRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		} else {
			resData.AnswerIndex.CollectCount += int32(collectCnt)
		}

		resData.AnswerIndex = answerIndexProto
	} else {
		logger.Errorf("get answerIndex cache failed, err: %v", err)

		answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex

		answerIndex, err := answerIndexModel.WithContext(l.ctx).
			Where(answerIndexModel.ID.Eq(in.AnswerId)).
			First()
		if err != nil {
			logger.Errorf("get answerIndex failed, err: mysql err, %v", err)
			res = &pb.GetAnswerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}

		answerIndexProto := &pb.AnswerIndex{
			Id:           answerIndex.ID,
			QuestionId:   answerIndex.QuestionID,
			UserId:       answerIndex.UserID,
			IpLoc:        answerIndex.IPLoc,
			ApproveCount: answerIndex.ApproveCount,
			LikeCount:    answerIndex.LikeCount,
			CollectCount: answerIndex.CollectCount,
			State:        answerIndex.State,
			Attrs:        answerIndex.Attrs,
			CreateTime:   answerIndex.CreateTime.String(),
			UpdateTime:   answerIndex.UpdateTime.String(),
		}

		resData.AnswerIndex = answerIndexProto

		answerIndexBytes, err = proto.Marshal(answerIndexProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v")
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("answer_index_%d", answerIndex.ID),
				answerIndexBytes,
				time.Second*86400)
		}
	}

	answerContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("answer_content_%d", in.AnswerId)).Bytes()
	if err == nil {
		answerContentProto := &pb.AnswerContent{}
		err = proto.Unmarshal(answerContentBytes, answerContentProto)
		if err != nil {
			logger.Errorf("unmarshal [answerContentProto] failed, err: %v", err)
			res = &pb.GetAnswerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
		resData.AnswerContent = answerContentProto
	} else {
		logger.Errorf("get answerContent cache failed, err: %v")

		answerContentModel := l.svcCtx.QuestionModel.AnswerContent

		answerContent, err := answerContentModel.WithContext(l.ctx).
			Where(answerContentModel.AnswerID.Eq(in.AnswerId)).
			First()
		if err != nil {
			logger.Errorf("get answerContent failed, err: mysql err, %v", err)
			res = &pb.GetAnswerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		answerContentProto := &pb.AnswerContent{
			AnswerId:   answerContent.AnswerID,
			Content:    answerContent.Content,
			Meta:       answerContent.Meta,
			CreateTime: answerContent.CreateTime.String(),
			UpdateTime: answerContent.UpdateTime.String(),
		}

		resData.AnswerContent = answerContentProto

		answerContentBytes, err = proto.Marshal(answerContentProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("answer_content_%d", answerContent.AnswerID),
				answerContentBytes,
				time.Second*86400)
		}
	}

	res = &pb.GetAnswerRes{
		Code: http.StatusOK,
		Msg:  "get answer successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
