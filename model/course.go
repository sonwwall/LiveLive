package model

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	Classname   string `gorm:"not null" json:"classname,required" form:"classname,required"`
	Description string
	TeacherId   int `gorm:"not null"`
}

type CourseMember struct {
	gorm.Model
	CourseId  int    `gorm:"not null" `
	StudentId int    `gorm:"not null"`
	Classname string `gorm:"not null" json:"classname,required" form:"classname,required"`
	JoinedAt  time.Time
}

type JoinCourse struct {
	StudentId      int    `gorm:"not null"`
	Classname      string ` json:"classname,required" form:"classname,required"`
	InvitationCode string ` json:"invitation_code,required" form:"invitation_code,required"`
}
