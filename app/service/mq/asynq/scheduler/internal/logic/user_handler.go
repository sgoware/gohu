package logic

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
)

func (l *AsyNqScheduler) updateUserRecord() {
	logger := log.GetSugaredLogger()

	payload, err := json.Marshal(job.ScheduleUpdateUserRecordPayload{})
	if err != nil {
		logger.Errorf("marshal ScheduleUpdateUserRecordPayload to json failed, err: %v", err)
	}

	task := asynq.NewTask(job.ScheduleUpdateUserRecord, payload)
	entryId, err := l.svcCtx.Scheduler.Register("*/1 * * * *", task)
	if err != nil {
		logger.Errorf("register [updateUserRecord] task to scheduler failed, err: %v", err)
	}
	logger.Debugf("scheduler registered an entry: %v", entryId)
}
