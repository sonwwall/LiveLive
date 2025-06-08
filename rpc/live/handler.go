package main

import (
	"LiveLive/dao/db"
	dao "LiveLive/dao/rdb"
	"LiveLive/kitex_gen/livelive/ai"
	"LiveLive/kitex_gen/livelive/ai/aiservice"
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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LiveServiceImpl implements the last service interface defined in the IDL.
type LiveServiceImpl struct {
	wsClient websocketservice.Client
	aiClient aiservice.Client
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
	pipe := dao.Redis.Pipeline() // 批量写入 Redis，提高性能
	key := fmt.Sprintf("sign:course:%d", course.ID)
	for _, member := range *students {

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

	s.wsClient.BroadcastToCourse(ctx, &websocket.BroadcastToCourseReq{
		CourseId: int64(course.ID),
		Data:     data,
	})

	//设置一个过期时间，统计完后自动过期
	duration := 3 * time.Minute

	//设置定时任务
	time.AfterFunc(time.Duration(req.Deadline)*time.Second, func() {
		s.wsClient.CountRegister(ctx, &websocket.CountRegisterReq{
			CourseId: int64(course.ID),
		})
		dao.Redis.Expire(ctx, key, duration) //统计完三分钟后过期
	})

	res := &live.PublishRegisterResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil
}

// StartRecording implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) StartRecording(ctx context.Context, req *live.StartRecordingReq) (resp *live.StartRecordingResp, err error) {
	existCourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.StartRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "课程不存在",
			},
		}
		return res, nil

	}
	if err != nil {
		res := &live.StartRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	RtmpUrl := fmt.Sprintf("rtmp://localhost:1935/live/")
	StreamKeyRow := fmt.Sprintf("teacher_%d_course_%d_%s", req.TeacherId, existCourse.ID, req.Classname)
	// 确保 recordings 目录存在
	_ = os.MkdirAll("./recordings", os.ModePerm)

	// 录制文件保存路径
	output := fmt.Sprintf("./recordings/%s_%d.flv", StreamKeyRow, time.Now().Unix())

	// 启动 ffmpeg
	//cmd := exec.Command("ffmpeg", "-i", RtmpUrl+StreamKeyRow, "-c", "copy", "-f", "flv", output)

	cmd := exec.Command("/opt/homebrew/bin/ffmpeg", "-i", RtmpUrl+StreamKeyRow, "-c", "copy", "-f", "flv", output)

	err = cmd.Start()
	if err != nil {
		res := &live.StartRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrStartRecording,
				Msg:  "录制失败:" + err.Error(),
			},
		}
		return res, nil
	}
	recordingMap[StreamKeyRow] = cmd

	go func() {
		err = cmd.Wait()
		delete(recordingMap, StreamKeyRow)

		if err != nil {
			fmt.Println("录制进程异常退出:", err)
		}

		fmt.Println("录制完成，开始转码为 wav")
		wavPath, err := ConvertFLVToWAV(output)
		if err != nil {
			fmt.Println("音频转码失败:", err)
		} else {
			fmt.Println("音频转码成功，路径:", wavPath)
		}

		result, _ := s.aiClient.AnalyzeAudio(ctx, &ai.AnalyzeAudioReq{
			WavPath: wavPath,
		})
		if result == nil {
			log.Println("内部错误")
		}
		if result.BaseResp.Code != 0 {
			log.Println("分析音频错误:", result.BaseResp.Msg)
		}
		log.Println("总结:", result.Summary)

		msg := map[string]interface{}{
			"type": "summary",
			"data": map[string]interface{}{
				"summary": result.Summary,
			},
		}

		data, _ := json.Marshal(msg)

		s.wsClient.BroadcastToCourse(ctx, &websocket.BroadcastToCourseReq{
			CourseId: int64(existCourse.ID),
			Data:     data,
		})
	}()

	res := &live.StartRecordingResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil
}

// ConvertFLVToWAV 提取 flv 中的音频并转换为 16kHz 单声道 wav 格式
func ConvertFLVToWAV(flvPath string) (string, error) {
	// 确保输出目录存在
	outputDir := "./audios"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("创建音频目录失败: %v", err)
	}

	// 生成 wav 输出路径
	base := strings.TrimSuffix(filepath.Base(flvPath), filepath.Ext(flvPath))
	wavPath := filepath.Join(outputDir, fmt.Sprintf("%s_%d.wav", base, time.Now().Unix()))

	// 构造 ffmpeg 命令
	//cmd := exec.Command("ffmpeg",
	//	"-i", flvPath,
	//	"-vn",                  // 不要视频
	//	"-acodec", "pcm_s16le", // 编码为 PCM
	//	"-ar", "16000", // 采样率 16kHz
	//	"-ac", "1", // 单声道
	//	wavPath,
	//)

	cmd := exec.Command("/opt/homebrew/bin/ffmpeg",
		"-i", flvPath,
		"-vn",                  // 去除视频
		"-acodec", "pcm_s16le", // PCM 编码
		"-ar", "16000", // 采样率 16kHz
		"-ac", "1", // 单声道
		wavPath,
	)

	fmt.Println("执行 ffmpeg 转换:", cmd.String())

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg 执行失败: %v", err)
	}

	return wavPath, nil
}

// StopRecording implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) StopRecording(ctx context.Context, req *live.StopRecordingReq) (resp *live.StopRecordingResp, err error) {
	existCourse, err := db.FindCourseByClassnameAndTeacherId(req.Classname, req.TeacherId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := &live.StopRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrCourseNotExist,
				Msg:  "课程不存在",
			},
		}
		return res, nil

	}
	if err != nil {
		res := &live.StopRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrDB,
				Msg:  "数据库错误：" + err.Error(),
			},
		}
		return res, nil
	}

	StreamKeyRow := fmt.Sprintf("teacher_%d_course_%d_%s", req.TeacherId, existCourse.ID, req.Classname)
	// 录制文件保存路径
	//output := fmt.Sprintf("./recordings/%s_%d.flv", StreamKeyRow, time.Now().Unix())
	proc, ok := recordingMap[StreamKeyRow]
	if !ok {
		res := &live.StopRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrStopRecording,
				Msg:  "结束录制失败",
			},
		}
		return res, nil
	}

	// 杀掉进程
	err = proc.Process.Kill()
	if err != nil {
		res := &live.StopRecordingResp{
			BaseResp: &base.BaseResp{
				Code: code.ErrStopRecording,
				Msg:  "结束录制失败:" + err.Error(),
			},
		}
		return res, nil
	}

	res := &live.StopRecordingResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil
}
