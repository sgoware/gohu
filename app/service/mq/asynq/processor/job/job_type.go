package job

const (
	MsgCreateUserSubjectTask            = "msg:user_subject:create"
	MsgUpdateUserSubjectRecordTask      = "msg:user_subject_record:update"
	MsgUpdateUserSubjectRecordCacheTask = "msg:user_subject_cache:update"
	MsgAddUserSubjectRecordCacheTask    = "msg:user_subject_cache:add"
	ScheduleUpdateUserSubjectRecordTask = "schedule:user_subject_record:update"

	ScheduleUpdateQuestionRecordTask = "schedule:question:update"
	ScheduleUpdateCommentRecordTask  = "schedule:comment:record"
	DeferTask                        = "defer:task:xx"
)
