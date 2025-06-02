package main

import (
	"LiveLive/dao/db"
	dao "LiveLive/dao/rdb"
	"LiveLive/kitex_gen/livelive/base"
	live "LiveLive/kitex_gen/livelive/live"
	"LiveLive/kitex_gen/livelive/websocket"
	"LiveLive/kitex_gen/livelive/websocket/websocketservice"
	"LiveLive/model"
	"LiveLive/rpc/live/code"
	"LiveLive/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

// LiveServiceImpl implements the last service interface defined in the IDL.
type LiveServiceImpl struct {
	wsClient websocketservice.Client
}

// GetStreamKey implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq) (resp *live.GetStreamKeyResp, err error) {
	existCourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.GetStreamKeyResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "课程不存在",
			},
		}
		return res, nil

	}
	if err != nil {
		res := &live.GetStreamKeyResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	RtmpUrl := fmt.Sprintf("rtmp://localhost:1935/live")
	StreamKeyRow := fmt.Sprintf("teacher_%d_course_%d_%s", req.TeacherId, existCourse.ID, req.Classname)
	log.Println("测试4")
	log.Println(StreamKeyRow)
	//TODO 注意这里如果没有连上livego不会报错！！！检查怎么回事
	StreamKey, err := utils.GetStreamKey(StreamKeyRow)
	log.Println("测试5")
	log.Println(StreamKey)
	if err != nil {
		res := &live.GetStreamKeyResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrGetStreamKey,
				Msg:  err.Error(),
			},
		}
		return res, nil
	}
	log.Println("测试6")

	livesession := &model.LiveSession{
		ClassName: req.Classname,
		TeacherID: req.TeacherId,
		CourseID:  int64(existCourse.ID),
		StartTime: time.Now(),
		RtmpURL:   RtmpUrl,
		StreamKey: StreamKey,
	}
	log.Println("测试2")

	err = db.AddLive(livesession)
	if err != nil {
		res := &live.GetStreamKeyResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  err.Error(),
			},
		}
		return res, nil
	}

	log.Println("测试3")

	res := &live.GetStreamKeyResp{
		RtmpUrl:   RtmpUrl,
		StreamKey: StreamKey,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	log.Println("测试1")
	log.Println(RtmpUrl, StreamKey)

	return res, nil
}

// WatchLive implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) WatchLive(ctx context.Context, req *live.WatchLiveReq) (resp *live.WatchLiveResp, err error) {

	teacher, err := db.FindUserByUsername(req.TeacherName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.WatchLiveResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrTeacherNotExist,
				Msg:  "该老师不存在",
			},
		}
		return res, nil
	}
	if err != nil {
		res := &live.WatchLiveResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	course, err := db.FindCourseByClassnameAndTeacherId(req.Classname, int64(teacher.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.WatchLiveResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "该课程不存在",
			},
		}
		return res, nil
	}
	if err != nil {
		res := &live.WatchLiveResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	_, err = db.FindCourseMemberByCourseIdAndStudentId(int64(course.ID), req.StudentId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.WatchLiveResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrNoCoursePermission,
				Msg:  "您未加入该课程！",
			},
		}
		return res, nil
	}

	//baseURL := "localhost:8060"
	//StreamKeyRow := fmt.Sprintf("teacher_%d,course_%d_%s", req.TeacherId, existCourse.ID, req.Classname)
	uri := fmt.Sprintf("/live/teacher_%d_course_%d_%s.flv", teacher.ID, course.ID, req.Classname)
	//log.Println(uri)
	//uid := req.StudentId    // 学生的用户 ID
	//secret := "sonwwall"    // 与 Nginx 统一
	//ttl := 10 * time.Minute // 有效期 10 分钟
	//
	//playURL := utils.GeneratePlayURL(baseURL, uri, secret, uid, ttl)
	res := &live.WatchLiveResp{
		Addr: uri,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil

}

// PublishRegister 签到
// PublishRegister implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) PublishRegister(ctx context.Context, req *live.PublishRegisterReq) (resp *live.PublishRegisterResp, err error) {
	students, err := db.FindStudentByTeacherNameAndClassName(req.TeacherName, req.Classname)
	if err != nil {
		res := &live.PublishRegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误:" + err.Error(),
			},
		}
		return res, nil
	}
	pipe := dao.Redis.Pipeline() // 批量写入 Redis，提高性能
	for _, member := range *students {
		key := fmt.Sprintf("sign:course:%d", member.CourseId)
		field := member.StudentName
		pipe.HSet(context.Background(), key, field, 0)
	}
	//写入
	_, err = pipe.Exec(context.Background())
	if err != nil {
		res := &live.PublishRegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "Redis写入错误:" + err.Error(),
			},
		}
		return res, nil
	}
	msg := map[string]interface{}{
		"type": "register",
		"data": map[string]interface{}{
			"code": 1,
		},
	}
	data, _ := json.Marshal(msg)
	course, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)
	if err != nil {
		res := &live.PublishRegisterResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误::" + err.Error(),
			},
		}
		return res, nil
	}
	s.wsClient.BroadcastToCourse(ctx, &websocket.BroadcastToCourseReq{
		CourseId: int64(course.ID),
		Data:     data,
	})

	res := &live.PublishRegisterResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil
}
