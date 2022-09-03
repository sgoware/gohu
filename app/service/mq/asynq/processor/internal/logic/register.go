package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"main/app/service/mq/asynq/processor/internal/logic/comment"
	"main/app/service/mq/asynq/processor/internal/logic/question"
	"main/app/service/mq/asynq/processor/internal/logic/user"
	"main/app/service/mq/asynq/processor/internal/svc"
	"main/app/service/mq/asynq/processor/job"
)

type Processor struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessor(ctx context.Context, svcCtx *svc.ServiceContext) *Processor {
	return &Processor{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (p *Processor) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.Handle(job.MsgCreateUserSubjectTask, user.NewCreateUserSubjectRecordHandler(p.svcCtx.Config))
	mux.Handle(job.MsgUpdateUserSubjectRecordTask, user.NewUpdateUserSubjectRecordHandler(p.svcCtx.Config))
	mux.Handle(job.MsgUpdateUserSubjectCacheTask, user.NewUpdateUserSubjectCacheHandler(p.svcCtx.Config))
	mux.Handle(job.MsgAddUserSubjectCacheTask, user.NewMsgAddUserSubjectCacheHandler(p.svcCtx.Config))
	mux.Handle(job.ScheduleUpdateUserSubjectRecordTask, user.NewScheduleUpdateUserSubjectRecordHandler(p.svcCtx.Config))
	mux.Handle(job.ScheduleUpdateUserCollectRecordTask, user.NewScheduleUpdateUserCollectRecordHandler(p.svcCtx.Config))

	mux.Handle(job.MsgCrudQuestionSubjectRecordTask, question.NewMsgCrudQuestionSubjectHandler(p.svcCtx.Config))
	mux.Handle(job.MsgCrudQuestionContentRecordTask, question.NewMsgCrudQuestionContentHandler(p.svcCtx.Config))
	mux.Handle(job.ScheduleUpdateQuestionSubjectRecordTask, question.NewScheduleUpdateQuestionSubjectRecordHandler(p.svcCtx.Config))
	mux.Handle(job.ScheduleUpdateAnswerIndexRecordTask, question.NewScheduleUpdateAnswerIndexRecordHandler(p.svcCtx.Config))

	mux.Handle(job.MsgCrudCommentSubjectTask, comment.NewMsgCrudCommentSubjectHandler(p.svcCtx.Config))

	return mux
}
