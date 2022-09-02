package svc

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/scheduler/internal/config"
	"time"
)

type ServiceContext struct {
	Config config.Config

	Scheduler *asynq.Scheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Scheduler: newScheduler(c),
	}
}

func newScheduler(c config.Config) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		c.RedisConf,
		&asynq.SchedulerOpts{
			Location: location,
			EnqueueErrorHandler: func(task *asynq.Task, opts []asynq.Option, err error) {
				logger := log.GetSugaredLogger()
				logger.Errorf("[Scheduler EnqueueErrorHandler] err: %+v, task: %+v", err, task)
			},
		})
}
