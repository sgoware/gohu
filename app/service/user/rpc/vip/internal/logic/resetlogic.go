package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/user/rpc/vip/internal/svc"
	"main/app/service/user/rpc/vip/pb"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetLogic) Reset(in *pb.ResetReq) (res *pb.ResetRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubjectModel := l.svcCtx.UserModel.UserSubject
	_, err = userSubjectModel.WithContext(l.ctx).
		Select(userSubjectModel.ID, userSubjectModel.Vip).
		Where(userSubjectModel.ID.Eq(in.Id)).
		Update(userSubjectModel.Vip, 0)
	if err != nil {
		logger.Errorf("reset vip failed, err: %v", err)
		res = &pb.ResetRes{
			Code: http.StatusOK,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	payload, err := json.Marshal(job.MsgUpdateUserSubjectCachePayload{
		Id:  in.Id,
		Vip: 0,
	})
	if err != nil {
		logger.Errorf("marshal [MsgUpdateUserSubjectCachePayload] to json failed, err: %v", err)
		res = &pb.ResetRes{
			Code: http.StatusOK,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgUpdateUserSubjectCacheTask, payload))
	if err != nil {
		logger.Errorf("create [MsgUpdateUserSubjectCacheTask] insert queue failed, err: %v", err)
		res = &pb.ResetRes{
			Code: http.StatusOK,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	res = &pb.ResetRes{
		Code: http.StatusOK,
		Msg:  "reset vip level successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
