package logic

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
)

func (l *AsyNqScheduler) updateCoomentRecord() {
	logger := log.GetSugaredLogger()

	entryId, err := l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateCommentSubjectTask, nil))
	if err != nil {
		logger.Errorf("register [ScheduleUpdateCommentSubjectTask] task to scheduler failed, err: %v", err)
	}

	logger.Debugf("scheduler registered an entry: %v", entryId)

	entryId, err = l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateCommentIndexTask, nil))
	if err != nil {
		logger.Errorf("register [ScheduleUpdateCommentIndexTask] task to scheduler failed, err: %v", err)
	}

	logger.Debugf("scheduler registered an entry: %v", entryId)
}
