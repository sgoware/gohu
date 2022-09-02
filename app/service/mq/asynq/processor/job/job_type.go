package job

const (
	MsgCreateUserSubjectTask            = "msg:user_subject:create"
	MsgUpdateUserSubjectRecordTask      = "msg:user_subject_record:update"
	MsgUpdateUserSubjectCacheTask       = "msg:user_subject_cache:update"
	MsgAddUserSubjectCacheTask          = "msg:user_subject_cache:add"
	ScheduleUpdateUserSubjectRecordTask = "schedule:user_subject_record:update"

	MsgCrudQuestionSubjectRecordTask = "msg:question_subject_record:crud"
	MsgCrudQuestionContentRecordTask = "msg:question_content_record:crud"

	MsgUpdateUserCollectCacheTask       = "msg:user_collect_cache:update"
	ScheduleUpdateUserCollectRecordTask = "schedule:user_collect_record:update"

	MsgCrudCommentSubjectTask = "msg:comment_subject:crud"

	ScheduleUpdateQuestionRecordTask = "schedule:question:update"
	ScheduleUpdateCommentRecordTask  = "schedule:comment:record"
	DeferTask                        = "defer:task:xx"
)
