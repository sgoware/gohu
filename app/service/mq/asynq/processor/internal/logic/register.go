package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"main/app/service/mq/asynq/processor/internal/config"
	"main/app/service/mq/asynq/processor/internal/logic/user"
	"main/app/service/mq/asynq/processor/internal/svc"
	"main/app/service/mq/asynq/processor/job"
)

type Processor struct {
	ctx    context.Context
	config config.Config
	svcCtx *svc.ServiceContext
}

func NewProcessor(ctx context.Context, c config.Config, svcCtx *svc.ServiceContext) *Processor {
	return &Processor{
		ctx:    ctx,
		config: c,
		svcCtx: svcCtx,
	}
}

func (p *Processor) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.Handle(job.MsgCreateUserSubjectTask, user.NewCreateUserSubjectRecordHandler(p.config))
	mux.Handle(job.MsgUpdateUserSubjectRecordTask, user.NewUpdateUserSubjectRecordHandler(p.config))

	return mux
}
