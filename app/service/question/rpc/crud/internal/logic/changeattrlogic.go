package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAttrLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAttrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAttrLogic {
	return &ChangeAttrLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeAttrLogic) ChangeAttr(in *pb.ChangeAttrReq) (res *pb.ChangeAttrRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	switch in.ObjType {
	case 1:
		// 问题
		questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject

		questionSubject, err := questionSubjectModel.WithContext(l.ctx).
			Where(questionSubjectModel.ID.Eq(in.ObjId)).
			First()
		if err != nil {
			logger.Errorf("change question attr failed, err: mysql err, %v", err)
			res = &pb.ChangeAttrRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		switch in.AttrType {
		case 1:
			// 关注
			switch in.Action {
			case 0:
				// 增加
				_, err = questionSubjectModel.WithContext(l.ctx).
					Where(questionSubjectModel.ID.Eq(in.ObjId)).
					Update(questionSubjectModel.SubCount, questionSubject.SubCount+1)
				if err != nil {
					logger.Errorf("increase question subscription failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}
			case 1:
				// 减少
				_, err = questionSubjectModel.WithContext(l.ctx).
					Where(questionSubjectModel.ID.Eq(in.ObjId)).
					Update(questionSubjectModel.SubCount, questionSubject.SubCount-1)
				if err != nil {
					logger.Errorf("decrease question subscription failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}
			}

		case 2:
			// 浏览
			// TODO: 问题浏览数
		}

	case 2:
		// 回答
		answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex

		answerIndex, err := answerIndexModel.WithContext(l.ctx).
			Where(answerIndexModel.ID.Eq(in.ObjId)).
			First()
		if err != nil {
			logger.Errorf("change answer attr failed, err: mysql err, %v", err)
			res = &pb.ChangeAttrRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		switch in.AttrType {
		case 0:
			// 赞同
			switch in.Action {
			case 0:
				// 增加
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.ApproveCount, answerIndex.ApproveCount+1)
				if err != nil {
					logger.Errorf("increase answer approve failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}

			case 1:
				// 减少
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.ApproveCount, answerIndex.ApproveCount-1)
				if err != nil {
					logger.Errorf("decrease answer approve failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}
			}

		case 1:
			// 喜欢
			switch in.Action {
			case 0:
				// 增加
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.LikeCount, answerIndex.LikeCount+1)
				if err != nil {
					logger.Errorf("increase answer like failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}

			case 1:
				// 减少
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.LikeCount, answerIndex.LikeCount-1)
				if err != nil {
					logger.Errorf("decrease answer like failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}
			}

		case 2:
			// 收藏
			switch in.Action {
			case 0:
				// 增加
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.CollectCount, answerIndex.CollectCount+1)
				if err != nil {
					logger.Errorf("increase answer collect failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}

			case 1:
				// 减少
				_, err = answerIndexModel.WithContext(l.ctx).
					Where(answerIndexModel.ID.Eq(in.ObjId)).
					Update(answerIndexModel.CollectCount, answerIndex.CollectCount-1)
				if err != nil {
					logger.Errorf("decrease answer collect failed, err: mysql err, %v", err)
					res = &pb.ChangeAttrRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					logger.Debugf("send message: %v", res.String())
					return res, nil
				}
			}
		}
	}

	res = &pb.ChangeAttrRes{
		Code: http.StatusOK,
		Msg:  "change attr successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
