package main

import (
	"LiveLive/dao"
	"LiveLive/dao/db"
	"LiveLive/model"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"time"
)

//定时任务，定时清理过期的邀请码

func StartInviteCodeCleaner(db *gorm.DB) {
	c := cron.New()

	// 每1分钟执行
	c.AddFunc("@every 1m", func() {
		fmt.Println("开始清理过期邀请码...")
		now := time.Now()
		result := db.Where("expired_at < ?", now).Delete(&model.CourseInvite{})
		if result.Error != nil {
			log.Println("删除失败:", result.Error)
		} else {
			log.Printf("成功清理 %d 个过期邀请码\n", result.RowsAffected)
		}
	})

	c.Start()
}

func main() {
	dao.Init() // 你的 GORM 初始化函数
	StartInviteCodeCleaner(db.Mysql)

	// 保持服务运行
	select {}
}
