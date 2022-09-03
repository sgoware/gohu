package logic

import (
	"context"
	"main/app/service/mq/asynq/scheduler/internal/svc"
)

type AsyNqScheduler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduler(ctx context.Context, svcCtx *svc.ServiceContext) *AsyNqScheduler {
	return &AsyNqScheduler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AsyNqScheduler) Register() {
	l.updateUserRecord()
	l.updateQuestionRecord()
}
