package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"main/app/service/mq/asynq/processor/internal/svc"
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

	//mux.Handle()

	return mux
}
