package db

import (
	"LiveLive/model"
	"gorm.io/gorm"
)

var Mysql *gorm.DB

func AddUser(user *model.User) *gorm.DB {
	return Mysql.Create(user)
}

func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := Mysql.Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserByMobile(mobile string) (*model.User, error) {
	var user model.User
	err := Mysql.Where("mobile = ?", mobile).First(&user).Error
	return &user, err
}

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := Mysql.Where("email = ?", email).First(&user).Error
	return &user, err
}

func CheckUser(account, password string) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := Mysql.Where(Mysql.Or("username = ?", account)).
		Where("password = ?", password).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil

}
