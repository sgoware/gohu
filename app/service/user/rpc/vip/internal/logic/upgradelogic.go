package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	"net/http"

	"main/app/service/user/rpc/vip/internal/svc"
	"main/app/service/user/rpc/vip/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeLogic {
	return &UpgradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpgradeLogic) Upgrade(in *pb.UpgradeReq) (res *pb.UpgradeRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubjectModel := l.svcCtx.UserModel.UserSubject
	userInfo, _ := userSubjectModel.WithContext(l.ctx).
		Select(userSubjectModel.ID, userSubjectModel.Vip).
		Where(userSubjectModel.ID.Eq(cast.ToInt64(in.Id))).First()
	if userInfo.Vip < 9 {
		_, err = userSubjectModel.WithContext(l.ctx).
			Select(userSubjectModel.ID, userSubjectModel.Vip).
			Where(userSubjectModel.ID.Eq(userInfo.ID)).
			Update(userSubjectModel.Vip, userInfo.Vip+1)
		if err != nil {
			logger.Errorf("update [user_suject] record failed, err: %v", err)
			res = &pb.UpgradeRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		payload, err := json.Marshal(job.MsgAddUserSubjectCachePayload{
			Id:  in.Id,
			Vip: userInfo.Vip + 1,
		})
		if err != nil {
			logger.Errorf("marshal [MsgAddUserSubjectCachePayload] to json failed, err: %v", err)
			res = &pb.UpgradeRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgAddUserSubjectCacheTask, payload))
		if err != nil {
			logger.Errorf("create MsgUpdateUserSubjectCacheTask] insert queue failed, err: %v", err)
			res = &pb.UpgradeRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		res = &pb.UpgradeRes{
			Code: http.StatusOK,
			Msg:  "vip upgrade successfully",
			Ok:   true,
			Data: &pb.UpgradeRes_Data{VipLevel: userInfo.Vip + 1},
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	} else {
		res = &pb.UpgradeRes{
			Code: http.StatusBadRequest,
			Msg:  "vip level is already the highest",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
}
