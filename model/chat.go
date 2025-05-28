package model

import "gorm.io/gorm"

type ChatMsgRecord struct {
	gorm.Model
	UserID   int64
	CourseID int64
	Content  string
}
