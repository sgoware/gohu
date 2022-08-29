// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameQuestionSubject = "question_subject"

// QuestionSubject mapped from table <question_subject>
type QuestionSubject struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`             // 主键
	UserID      int64     `gorm:"column:user_id;not null" json:"user_id"`                        // 提问者 id
	IPLoc       string    `gorm:"column:ip_loc;not null" json:"ip_loc"`                          // 提问者 IP 归属地
	Title       string    `gorm:"column:title;not null" json:"title"`                            // 问题标题
	Topic       string    `gorm:"column:topic" json:"topic"`                                     // 问题主题
	Tag         string    `gorm:"column:tag" json:"tag"`                                         // 问题标签
	SubCount    int32     `gorm:"column:sub_count;not null" json:"sub_count"`                    // 关注总数
	AnswerCount int32     `gorm:"column:answer_count;not null" json:"answer_count"`              // 回答总数
	ViewCount   int64     `gorm:"column:view_count;not null" json:"view_count"`                  // 浏览总数
	State       int32     `gorm:"column:state;not null" json:"state"`                            // 状态 (0-正常 1-隐藏)
	Attrs       int32     `gorm:"column:attrs" json:"attrs"`                                     // 属性 (待添加)
	CreateTime  time.Time `gorm:"autoCreateTime;column:create_time;not null" json:"create_time"` // 创建时间
	UpdateTime  time.Time `gorm:"autoUpdateTime;column:update_time;not null" json:"update_time"` // 修改时间
}

// TableName QuestionSubject's table name
func (*QuestionSubject) TableName() string {
	return TableNameQuestionSubject
}
