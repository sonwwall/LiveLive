package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	user "LiveLive/kitex_gen/livelive/user"
	"LiveLive/model"
	"LiveLive/rpc/user/code"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	usr := &model.User{Username: req.Username, Password: req.Password, Mobile: req.Mobile, Email: req.Email, Role: req.Role}
	log.Println("role1", usr.Role)
	//进行参数验证
	existUser, err := db.FindUserByUsername(req.Username)
	if err == nil || existUser == nil {
		res := &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrUsernameExist,
				Msg:  "该用户名已存在",
			},
		}
		return res, nil

	}

	existUser, err = db.FindUserByMobile(req.Mobile)
	if err == nil || existUser == nil {
		res := &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrPhoneExist,
				Msg:  "该手机号已经注册过",
			},
		}
		return res, nil
	}

	existUser, err = db.FindUserByEmail(req.Email)
	if err == nil || existUser == nil {
		res := &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrEmailExist,
				Msg:  "该邮箱已被注册过",
			},
		}
		return res, nil
	}

	result := db.AddUser(usr)
	if result.Error != nil {
		res := &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + result.Error.Error(),
			},
		}

		return res, nil
	}
	res := &user.RegisterResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}

	return res, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	usr := &model.User{Username: req.Username, Password: req.Password}

	result, resultErr := db.FindUserByUsername(usr.Username)
	if resultErr != nil || result == nil {
		if errors.Is(resultErr, gorm.ErrRecordNotFound) {
			res := &user.LoginResp{

				BaseResp: &base.BaseResp{
					Code: code.ErrUserNotExists,
					Msg:  "用户不存在，请先注册",
				},
			}
			return res, nil
		}
		res := &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: 2003,
				Msg:  resultErr.Error(),
			},
		}
		return res, nil
	}

	if result.Password != req.Password {
		res := &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: 2004,
				Msg:  "密码错误",
			},
		}
		return res, nil
	}

	res := &user.LoginResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}

	return res, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoReq) (resp *user.UserInfoResp, err error) {

	username := req.GetUsername()
	result, resultErr := db.FindUserByUsername(username)
	if resultErr != nil || result == nil {
		if errors.Is(resultErr, gorm.ErrRecordNotFound) {
			res := &user.UserInfoResp{

				BaseResp: &base.BaseResp{
					Code: code.ErrUserNotExists,
					Msg:  "用户不存在",
				},
			}
			return res, nil
		}
		res := &user.UserInfoResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + resultErr.Error(),
			},
		}
		return res, nil
	}

	res := &user.UserInfoResp{
		Username: result.Username,
		Email:    result.Email,
		Mobile:   result.Mobile,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}

	return res, nil
}
