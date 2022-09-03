package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/question/rpc/info/internal/svc"
	"main/app/service/question/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionLogic) GetQuestion(in *pb.GetQuestionReq) (res *pb.GetQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetQuestionRes_Data{}

	questionSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("question_subject_%d", in.QuestionId)).Bytes()
	if err != nil {
		logger.Errorf("get questionSubject cache failed, err: %v", err)

		questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject

		questionSubject, err := questionSubjectModel.WithContext(l.ctx).
			Where(questionSubjectModel.ID.Eq(in.QuestionId)).
			First()
		if err != nil {
			logger.Errorf("get questionSubject failed, err: mysql err, %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		questionSubjectProto := &pb.QuestionSubject{
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
			Attr:        questionSubject.Attrs,
			CreateTime:  questionSubject.CreateTime.String(),
			UpdateTime:  questionSubject.UpdateTime.String(),
		}

		resData.QuestionSubject = questionSubjectProto

		questionSubjectBytes, err = proto.Marshal(questionSubjectProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("question_subject_%d", questionSubject.ID),
				questionSubjectBytes,
				time.Second*86400)
		}
	} else {
		questionSubjectProto := &pb.QuestionSubject{}
		err = proto.Unmarshal(questionSubjectBytes, questionSubjectProto)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		resData.QuestionSubject = questionSubjectProto

		subCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("question_subject_sub_cnt_%d", in.QuestionId)).Int()
		if err != nil {
			logger.Errorf("get [question_subject_sub_cnt] failed, err: %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		answerCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("question_subject_answer_cnt_%d", in.QuestionId)).Int()
		if err != nil {
			logger.Errorf("get [question_subject_answer_cnt] failed, err: %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		viewCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("question_subject_view_cnt_%d", in.QuestionId)).Int64()
		if err != nil {
			logger.Errorf("get [question_subject_view_cnt] failed, err: %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		resData.QuestionSubject.SubCount += int32(subCnt)
		resData.QuestionSubject.AnswerCount += int32(answerCnt)
		resData.QuestionSubject.ViewCount += viewCnt
	}

	questionContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("question_content_%d", in.QuestionId)).Bytes()
	if err != nil {
		logger.Errorf("get questionContent cache failed, err: %v", err)

		questionContentModel := l.svcCtx.QuestionModel.QuestionContent

		questionContent, err := questionContentModel.WithContext(l.ctx).
			Where(questionContentModel.QuestionID.Eq(in.QuestionId)).
			First()
		if err != nil {
			logger.Errorf("get questionContent failed, err: mysql err, %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		questionContentProto := &pb.QuestionContent{
			QuestionId: questionContent.QuestionID,
			Content:    questionContent.Content,
			Meta:       questionContent.Meta,
			CreateTime: questionContent.CreateTime.String(),
			UpdateTime: questionContent.UpdateTime.String(),
		}

		resData.QuestionContent = questionContentProto

		questionContentBytes, err = proto.Marshal(questionContentProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("question_content_%d", questionContent.QuestionID),
				questionContentBytes,
				time.Second*86400)
		}
	} else {
		questionContentProto := &pb.QuestionContent{}
		err = proto.Unmarshal(questionContentBytes, questionContentProto)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetQuestionRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		resData.QuestionContent = questionContentProto
	}

	err = l.svcCtx.Rdb.Incr(l.ctx,
		fmt.Sprintf("question_subject_view_cnt_%d", in.QuestionId)).Err()
	if err != nil {
		logger.Errorf("incr [question_subject_view_cnt] failed, err: %v", err)
	} else {
		err = l.svcCtx.Rdb.SAdd(l.ctx,
			"question_subject_view_cnt_set",
			in.QuestionId).Err()
		if err != nil {
			logger.Errorf("update [question_subject_view_cnt_set] failed, err: %v", err)
		}
	}

	res = &pb.GetQuestionRes{
		Code: http.StatusOK,
		Msg:  "get question successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
