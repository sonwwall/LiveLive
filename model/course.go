package model

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	Classname   string `gorm:"not null"`
	Description string
	TeacherId   int `gorm:"not null"`
}

type CourseMember struct {
	gorm.Model
	CourseId  int `gorm:"not null"`
	StudentId int `gorm:"not null"`
	JoinedAt  time.Time
}
