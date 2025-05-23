package model

import "time"

// LiveSession 直播记录表
type LiveSession struct {
	ID        int64  `gorm:"primaryKey"`
	CourseID  int64  `gorm:"not null"` // 所属课程
	TeacherID int64  `gorm:"not null"` // 冗余字段，便于查询
	RtmpURL   string `gorm:"not null"`
	StreamKey string `gorm:"not null"` // RTMP 推流码
	ClassName string `gorm:"not null" json:"classname,required" form:"classname,required"`
	StartTime time.Time
	EndTime   *time.Time //可空时间，直播进行时时可以为nil
}
