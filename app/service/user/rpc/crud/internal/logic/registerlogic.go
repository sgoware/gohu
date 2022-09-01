package logic

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/sha3"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"main/app/utils/uuid"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (res *pb.RegisterRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if len(strings.TrimSpace(in.Username)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &crud.RegisterRes{
			Code: http.StatusBadRequest,
			Msg:  "param err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	ok, err := l.svcCtx.Rdb.SIsMember(l.ctx,
		"user_register",
		in.Username).Result()
	if err == nil {
		if !ok {
			logger.Debugf("[user_register] cache not found")
		} else {
			res = &crud.RegisterRes{
				Code: http.StatusForbidden,
				Msg:  "user already exist",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	} else {
		logger.Errorf("get [user_register] cache member failed, err: %v", err)
	}

	// 获取注册缓存失败, 看看数据库里有没有
	userSubjectModel := l.svcCtx.UserModel.UserSubject
	_, err = userSubjectModel.WithContext(l.ctx).Where(userSubjectModel.Username.Eq(in.Username)).First()
	switch err {
	case nil:
		// 用户已经存在的情况
		{
			res = &crud.RegisterRes{
				Code: http.StatusForbidden,
				Msg:  "user already exist",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	case gorm.ErrRecordNotFound:
		// 用户不存在的情况(可以创建用户)
		{
			// 密码使用sha3哈希然后存储
			d := sha3.Sum224([]byte(in.Password))
			encryptedPassword := hex.EncodeToString(d[:])

			// 生成默认昵称
			defaultNickname := fmt.Sprintf("gohu_%s", uuid.NewRandomString(in.Username, "username", 10))

			// 添加注册缓存
			err = l.svcCtx.Rdb.SAdd(l.ctx,
				"user_register",
				in.Username).Err()
			if err != nil {
				logger.Errorf("update user register cache failed, err: %v", err)
			}

			// 将注册信息加入消息队列(写入数据库)
			payload, err := json.Marshal(job.MsgCreateUserSubjectPayload{
				Username:   defaultNickname,
				Password:   encryptedPassword,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
			if err != nil {
				logger.Errorf("marshal [MsgUpdateUserSubjectRecordPayload] into json failed, err: %v", err)
			} else {
				_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgCreateUserSubjectTask, payload))
				if err != nil {
					logger.Errorf("create [MsgUpdateUserSubjectRecordTask] insert queue failed, err: %v", err)
				}
			}

			res = &crud.RegisterRes{
				Code: http.StatusOK,
				Msg:  "create user successfully",
				Ok:   true,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		// 数据库查询失败的情况
	default:
		logger.Errorf("query [user_subject] failed in mysql failed, err: %v", err)
		res = &crud.RegisterRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}
}
