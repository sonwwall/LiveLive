package model

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	Classname   string `gorm:"index:idx_teacher_classname,unique;type:varchar(100);not null" json:"classname,required" form:"classname,required"`
	Description string
	TeacherId   int `gorm:"index:idx_teacher_classname,unique;not null"` //这里使用了组合唯一约束
}

type CourseMember struct {
	gorm.Model

	CourseId  int `gorm:"not null"`
	StudentId int `gorm:"not null"`

	StudentName string `gorm:"not null;uniqueIndex:idx_class_teacher_student;type:varchar(100)"`

	TeacherName string `gorm:"not null;uniqueIndex:idx_class_teacher_student;type:varchar(100)"`

	Classname string `gorm:"not null;uniqueIndex:idx_class_teacher_student;type:varchar(100)" json:"classname,required" form:"classname,required"`

	JoinedAt time.Time
}

type JoinCourse struct {
	StudentId      int    `gorm:"not null"`
	Classname      string ` json:"classname,required" form:"classname,required"`
	InvitationCode string ` json:"invitation_code,required" form:"invitation_code,required"`
	TeacherName    string ` json:"teacher_name,required" form:"teacher_name,required"`
}
