package model

import "time"

type CourseInvite struct {
	ID         int64      `gorm:"primaryKey"`
	CourseID   int64      `gorm:"index"`
	Classname  string     `json:"classname,required" form:"classname,required"`
	Code       string     `gorm:"type:varchar(64);uniqueIndex"` // 邀请码本身
	MaxUsage   *int64     `json:"max_usage" form:"max_usage"`   // 可选，最大使用次数，=nil无限制，!=nil有限次
	UsageCount int64      // 已使用次数
	ExpiredAt  *time.Time `json:"expired_at" form:"expired_at"` // 过期时间，可选
	CreatedAt  time.Time
}
