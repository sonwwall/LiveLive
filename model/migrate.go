package model

import (
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"gorm.io/gorm"
)

func MigrateCourseAndCourseMember(db *gorm.DB) {
	err := db.AutoMigrate(&Course{}, &CourseMember{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}

func MigrateCourseInvite(db *gorm.DB) {
	err := db.AutoMigrate(&CourseInvite{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}

func MigrateUser(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}

func MigrateLive(db *gorm.DB) {
	err := db.AutoMigrate(&LiveSession{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}

func MigrateQuestion(db *gorm.DB) {
	err := db.AutoMigrate(&ChoiceQuestion{}, &AnswerChoiceQuestion{}, &AnsweredChoiceQuestion{})
	if err != nil {
		log.Errorf("迁移失败%s", err.Error())
	}
}

func MigrateChatMessage(db *gorm.DB) {
	err := db.AutoMigrate(&ChatMsgRecord{})
	if err != nil {
		log.Errorf("迁移失败%s", err.Error())
	}
}
