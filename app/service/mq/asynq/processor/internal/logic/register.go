package logic

import (
	"context"
	"github.com/hibiken/asynq"
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

	return mux
}
