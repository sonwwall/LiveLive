package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	course "LiveLive/kitex_gen/livelive/course"
	"LiveLive/model"
	"LiveLive/rpc/course/code"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct{}

// CreateCourse implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) CreateCourse(ctx context.Context, req *course.CreateCourseReq) (resp *course.CreateCourseResp, err error) {
	mycourse := &model.Course{Classname: req.Classname, Description: req.Description, TeacherId: int(req.TeacherId)}

	//参数验证
	existCourse, err := db.FindCourseByClassname(req.Classname)
	if err == nil && existCourse != nil {
		res := &course.CreateCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseExist,
				Msg:  "该课程已存在",
			},
		}
		return res, nil

	}

	err = db.CreateCourse(mycourse)
	if err != nil {
		res := &course.CreateCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误" + err.Error(),
			},
		}
		return res, nil
	}
	res := &course.CreateCourseResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil

}

// JoinCourse implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) JoinCourse(ctx context.Context, req *course.JoinCourseReq) (resp *course.JoinCourseResp, err error) {

	existcourse, err := db.FindCourseByClassname(req.Classname)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "课程不存在",
			},
		}
		return res, nil
	}
	if err != nil {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}
	joincourse := &model.CourseMember{Classname: req.Classname, StudentId: int(req.StudentId), CourseId: int(existcourse.ID), JoinedAt: time.Now()}
	err = db.AddStudentCourse(joincourse)
	if err != nil {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}
	res := &course.JoinCourseResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}

	return res, nil
}
