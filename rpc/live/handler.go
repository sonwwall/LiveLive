package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	live "LiveLive/kitex_gen/livelive/live"
	"LiveLive/model"
	"LiveLive/rpc/live/code"
	"LiveLive/utils"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

// LiveServiceImpl implements the last service interface defined in the IDL.
type LiveServiceImpl struct{}

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
	StreamKeyRow := fmt.Sprintf("teacher_%d,course_%d_%s", req.TeacherId, existCourse.ID, req.Classname)
	log.Println(StreamKeyRow)
	StreamKey, err := utils.GetStreamKey(StreamKeyRow)
	if err != nil {
		res := &live.GetStreamKeyResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  err.Error(),
			},
		}
		return res, nil
	}

	livesession := &model.LiveSession{
		ClassName: req.Classname,
		TeacherID: req.TeacherId,
		CourseID:  int64(existCourse.ID),
		StartTime: time.Now(),
		RtmpURL:   RtmpUrl,
		StreamKey: StreamKey,
	}

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

	res := &live.GetStreamKeyResp{
		RtmpUrl:   RtmpUrl,
		StreamKey: StreamKey,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}

	return res, nil
}
