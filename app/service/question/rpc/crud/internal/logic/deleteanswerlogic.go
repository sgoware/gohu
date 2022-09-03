package logic

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"main/app/common/log"
	"net/http"

	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAnswerLogic {
	return &DeleteAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAnswerLogic) DeleteAnswer(in *pb.DeleteAnswerReq) (res *pb.DeleteAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	err = l.svcCtx.Rdb.Del(l.ctx,
		fmt.Sprintf("answer_index_%d", in.AnswerId)).Err()
	if err != nil {
		logger.Errorf("delete [question_subject] cache failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Del(l.ctx,
		fmt.Sprintf("question_content_%d", in.AnswerId))
	if err != nil {
		logger.Errorf("delete [question_subject] cache failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 删除answer_index后, 级联删除关联的回答内容和评论
	answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex

	answerIndex, err := answerIndexModel.WithContext(l.ctx).
		Select(answerIndexModel.ID,
			answerIndexModel.QuestionID).
		Where(answerIndexModel.ID.Eq(in.AnswerId)).
		First()
	if err != nil {
		logger.Errorf("query [answer_index] record failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = l.svcCtx.Rdb.Decr(l.ctx,
		fmt.Sprintf("question_subject_answer_cnt_%d", answerIndex.QuestionID)).Err()
	if err != nil {
		logger.Errorf("decrease [question_subject_answer_cnt] failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	err = l.svcCtx.Rdb.SAdd(l.ctx,
		"question_subject_answer_cnt_set",
		answerIndex.QuestionID).Err()
	if err != nil {
		if err != nil {
			logger.Errorf("update [question_subject_sub_cnt_set] failed, err: %v", err)
			res = &pb.DeleteAnswerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}

	_, err = answerIndexModel.WithContext(l.ctx).Where(answerIndexModel.ID.Eq(in.AnswerId)).Delete()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.DeleteAnswerRes{
			Code: http.StatusBadRequest,
			Msg:  "answer not found",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	case nil:

	default:
		logger.Errorf("update Answer failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = l.svcCtx.Rdb.SRem(l.ctx,
		fmt.Sprintf("answer_id_user_set_%v", answerIndex.UserID),
		in.AnswerId).Err()
	if err != nil {
		logger.Errorf("update [answer_id_user_set] failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.DeleteAnswerRes{
		Code: http.StatusOK,
		Msg:  "delete answer successfully",
		Ok:   true,
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
