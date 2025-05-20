package model

import (
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username" form:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password" form:"password"`
}

func MigrateUser(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}
