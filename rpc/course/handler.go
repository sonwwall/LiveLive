package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	course "LiveLive/kitex_gen/livelive/course"
	"LiveLive/model"
	"LiveLive/rpc/course/code"
	"LiveLive/utils"
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
	existCourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)
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

	teacher, err := db.FindUserByUsername(req.TeacherName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrTeacherNotExist,
				Msg:  "该老师不存在！",
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

	existcourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, int64(teacher.ID))

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

	existCourseInvite, err := db.FindCourseInviteByCode(req.InvitationCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrInviteCodeInvalid,
				Msg:  "邀请码无效",
			},
		}
		return res, nil
	}

	if int64(existcourse.ID) != existCourseInvite.CourseID {
		res := &course.JoinCourseResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrInviteCodeInvalid,
				Msg:  "邀请码无效",
			},
		}
		return res, nil
	}

	joincourse := &model.CourseMember{
		Classname:   req.Classname,
		StudentId:   int(req.StudentId),
		CourseId:    int(existcourse.ID),
		JoinedAt:    time.Now(),
		StudentName: req.StudentName,
		TeacherName: req.TeacherName,
	}
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

// CreateCourseInvite implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) CreateCourseInvite(ctx context.Context, req *course.CreateCourseInviteReq) (resp *course.CreateCourseInviteResp, err error) {
	existcourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)

	//检查一下课程是否存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &course.CreateCourseInviteResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "课程不存在",
			},
		}
		return res, nil
	}
	if err != nil {
		res := &course.CreateCourseInviteResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	//检查一下是否有没过期的邀请码，有的话直接返回
	existCourseInvite, err := db.FindCourseInviteByCourseId(existcourse.ID)
	if err == nil && existCourseInvite != nil {
		res := &course.CreateCourseInviteResp{
			InviteCode: existCourseInvite.Code,
			BaseResp: &base.BaseResp{
				Code: 0,
				Msg:  "ok",
			},
		}
		return res, nil
	}

	var inviteCode string
	for {
		//检查一下code是否存在
		tryinviteCode := utils.GenerateInviteCode(6)
		_, err := db.FindCourseInviteByCode(tryinviteCode)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			inviteCode = tryinviteCode
			break
		}

	}

	courseInvite := &model.CourseInvite{
		Classname: req.Classname,
		CourseID:  int64(existcourse.ID),
		MaxUsage:  &req.MaxUsage,
		CreatedAt: time.Now(),
		ExpiredAt: utils.TimestampToPtr(req.ExpiredAt),
		Code:      inviteCode,
	}

	err = db.AddCourseInvite(courseInvite)

	res := &course.CreateCourseInviteResp{
		InviteCode: inviteCode,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil
}
