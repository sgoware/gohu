package logic

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
)

func (l *AsyNqScheduler) updateUserRecord() {
	logger := log.GetSugaredLogger()

	task := asynq.NewTask(job.ScheduleUpdateUserSubjectRecordTask, nil)
	entryId, err := l.svcCtx.Scheduler.Register("*/1 * * * *", task)
	if err != nil {
		logger.Errorf("register [updateUserRecord] task to scheduler failed, err: %v", err)
	}
	logger.Debugf("scheduler registered an entry: %v", entryId)
}
