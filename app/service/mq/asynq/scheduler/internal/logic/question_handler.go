package logic

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
)

func (l *AsyNqScheduler) updateQuestionRecord() {
	logger := log.GetSugaredLogger()

	entryId, err := l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateQuestionSubjectRecordTask, nil))
	if err != nil {
		logger.Errorf("register [ScheduleUpdateQuestionSubjectRecordTask] task to scheduler failed, err: %v", err)
	}

	logger.Debugf("scheduler registered an entry: %v", entryId)

	entryId, err = l.svcCtx.Scheduler.Register("*/1 * * * *",
		asynq.NewTask(job.ScheduleUpdateAnswerIndexRecordTask, nil))
	if err != nil {
		logger.Errorf("register [ScheduleUpdateAnswerIndexRecordTask] task to scheduler failed, err: %v", err)
	}

	logger.Debugf("scheduler registered an entry: %v", entryId)
}
