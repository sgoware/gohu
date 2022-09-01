package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func getBasicAuth(clientId, clientSecret, userId string) (basicAuthString string) {
	encodeAuthString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s:%s", clientId, clientSecret, userId)))
	basicAuthString = "Basic " + encodeAuthString
	return
}

func (l *LoginLogic) Login(in *pb.LoginReq) (res *pb.LoginRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())
	// 判断传入参数是否为空
	if len(strings.TrimSpace(in.Username)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &pb.LoginRes{
			Code: http.StatusBadRequest,
			Msg:  "param not fit",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 在缓存中查找用户
	userLoginCache, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("user_login_%s", in.Username)).Result()
	if err == nil {
		output := strings.Split(userLoginCache, ":")
		userId := output[0]
		password := output[1]
		if in.Password == password {
			// // 更新最近登录 ip
			payload, err := json.Marshal(job.MsgUpdateUserSubjectRecordPayload{
				Id:     cast.ToInt64(userId),
				LastIp: in.LastIp,
			})
			if err != nil {
				logger.Errorf("marshal [MsgUpdateUserSubjectRecordPayload] into json failed, err: %v", err)
			} else {
				_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgUpdateUserSubjectRecordTask, payload))
				if err != nil {
					logger.Errorf("create [MsgUpdateUserSubjectRecordTask] insert queue failed, err: %v", err)
				}
			}

			// 生成 oauth 服务器的认证头
			basicAuthString := getBasicAuth(l.svcCtx.ClientId, l.svcCtx.ClientSecret, cast.ToString(userId))
			res = &pb.LoginRes{
				Code: http.StatusOK,
				Msg:  "login successfully",
				Ok:   true,
				Data: &pb.LoginRes_Data{AuthToken: basicAuthString},
			}
		}
	}
	switch err {
	case redis.Nil:
		logger.Debugf("[user_login] cache not found")
	case nil:

	default:
		logger.Errorf("get [user_login] cache failed, err: %v", err)
	}

	// 在数据库中查找用户
	userSubjectModel := l.svcCtx.UserModel.UserSubject
	userInfo, err := userSubjectModel.WithContext(l.ctx).Where(userSubjectModel.Username.Eq(in.Username)).First()
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		// 设置空缓存,防止大量非法请求造成缓存穿透
		err = l.svcCtx.Rdb.Set(l.ctx,
			fmt.Sprintf("user_login_%s", in.Username),
			fmt.Sprintf("%d:%d", 0, 0),
			time.Second*86400).Err()
		if err != nil {
			logger.Errorf("update [user_login] cache failed, err: %v", err)
		}
		res = &pb.LoginRes{
			Code: http.StatusNotFound,
			Msg:  "uid not exist",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	default:
		{
			logger.Errorf("query [user_subject] in mysql failed, err: %v", err)
			res = &pb.LoginRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, err
		}
	}
	logger.Debugf("userInfo: \n%v", userInfo)

	// 验证密码
	d := sha3.Sum224([]byte(in.Password))
	encryptedPassword := hex.EncodeToString(d[:])
	if userInfo.Password != encryptedPassword {
		res = &pb.LoginRes{
			Code: http.StatusUnauthorized,
			Msg:  "wrong password",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 更新最近登录 ip
	_, err = userSubjectModel.WithContext(l.ctx).
		Where(userSubjectModel.Username.Eq(in.Username)).
		Update(userSubjectModel.LastIP, in.LastIp)
	if err != nil {
		logger.Errorf("database err, err: %v", err)
		res = &pb.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		return res, err
	}

	// 生成 oauth 服务器的认证头
	basicAuthString := getBasicAuth(l.svcCtx.ClientId, l.svcCtx.ClientSecret, cast.ToString(userInfo.ID))

	res = &pb.LoginRes{
		Code: http.StatusOK,
		Msg:  "get auth_token successfully",
		Ok:   true,
		Data: &pb.LoginRes_Data{AuthToken: basicAuthString},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
