package model

import (
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username" form:"username" binding:"required"`
	Password string `gorm:"type:varchar(255);not null" json:"password" form:"password" binding:"required"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" form:"email" binding:"required,email"`
	Mobile   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"mobile" form:"mobile" binding:"required"`
}

func MigrateUser(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}
