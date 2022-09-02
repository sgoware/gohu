package logic

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
)

func (l *AsyNqScheduler) updateUserRecord() {
	logger := log.GetSugaredLogger()

	entryId, err := l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateUserSubjectRecordTask, nil))
	if err != nil {
		logger.Errorf("register [updateUserSubjectRecord] task to scheduler failed, err: %v", err)
	}
	logger.Debugf("scheduler registered an entry: %v", entryId)

	entryId, err = l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateUserCollectRecordTask, nil))
	if err != nil {
		logger.Errorf("register [updateUserCollectRecord] task to scheduler failed, err: %v", err)
	}

	logger.Debugf("scheduler registered an entry: %v", entryId)
}
