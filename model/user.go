package model

import (
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username,required" form:"username,required" binding:"required" vd:"regexp('^\\w+$') && len($)>0"` //非空，字母数字或下划线
	Password string `gorm:"type:varchar(255);not null" json:"password,required" form:"password,required" binding:"required" vd:"len($)>=6"`                                //密码长度大于等于6
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email,required" form:"email,required" binding:"required,email" vd:"email($)"`                     //符合email格式
	Mobile   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"mobile,required" form:"mobile,required" binding:"required" vd:"phone($,'CN')"`                    //符合中国手机号格式
}

func MigrateUser(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Errorf("迁移失败：%s", err.Error())
	}
}
