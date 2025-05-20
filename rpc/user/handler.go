package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	user "LiveLive/kitex_gen/livelive/user"
	"LiveLive/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	usr := &model.User{Username: req.Username, Password: req.Password, Mobile: req.Mobile, Email: req.Email}
	result := db.AddUser(usr)
	if result.Error != nil {
		res := &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: "2001",
				Msg:  result.Error.Error(),
			},
		}

		return res, nil
	}
	res := &user.RegisterResp{
		BaseResp: &base.BaseResp{
			Code: "0",
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
					Code: "2002",
					Msg:  "用户不存在，请先注册",
				},
			}
			return res, nil
		}
		res := &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: "2003",
				Msg:  resultErr.Error(),
			},
		}
		return res, nil
	}

	if result.Password != req.Password {
		res := &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: "2004",
				Msg:  "密码错误",
			},
		}
		return res, nil
	}

	res := &user.LoginResp{
		BaseResp: &base.BaseResp{
			Code: "0",
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
					Code: "2002",
					Msg:  "用户不存在",
				},
			}
			return res, nil
		}
		res := &user.UserInfoResp{
			BaseResp: &base.BaseResp{
				Code: "2003",
				Msg:  resultErr.Error(),
			},
		}
		return res, nil
	}

	res := &user.UserInfoResp{
		Username: result.Username,
		Email:    result.Email,
		Mobile:   result.Mobile,
		BaseResp: &base.BaseResp{
			Code: "0",
			Msg:  "ok",
		},
	}

	return res, nil
}
