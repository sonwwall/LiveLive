package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ChoiceQuestion struct {
	gorm.Model
	CourseID  int64          `gorm:"not null"`
	TeacherId int64          `gorm:"not null"`
	Title     string         `gorm:"not null"`
	Options   datatypes.JSON `gorm:"not null"`
	Answer    int8           `gorm:"not null"`
	Deadline  *time.Time
}

type TrueOrFalseQuestion struct {
	gorm.Model
	TeacherId int64  `gorm:"not null"`
	CourseId  int64  `gorm:"not null"`
	Title     string `gorm:"not null"`
	Answer    int8   `gorm:"not null"`
	Deadline  *time.Time
}

type AnswerChoiceQuestion struct {
	gorm.Model
	ChoiceQuestionId int64 `gorm:"not null;uniqueIndex:uq_answer_question"` //一个学生只能答一次题
	StudentID        int64 `gorm:"not null;uniqueIndex:uq_answer_question"`
	Answer           *int
}

type AnswerTrueOrFalseQuestion struct {
	gorm.Model
	TrueOrFalseQuestionId int64 `gorm:"not null;uniqueIndex:uq_answer_question"`
	StudentID             int64 `gorm:"not null;uniqueIndex:uq_answer_question"`
	Answer                *int
}

type AnsweredChoiceQuestion struct {
	gorm.Model
	ChoiceQuestionId uint           `gorm:"not null;uniqueIndex:uq_answer_question"`
	Title            string         `gorm:"not null"`
	Options          datatypes.JSON `gorm:"not null"`
	Answer           int8
	Accuracy         float64
}

type AnsweredTrueOrFalseQuestion struct {
	gorm.Model
	TrueOrFalseQuestionId uint   `gorm:"not null;uniqueIndex:uq_answer_question"`
	Title                 string `gorm:"not null"`
	Answer                int8
	Accuracy              float64
}
