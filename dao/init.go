package dao

import (
	"LiveLive/dao/db"
	dao "LiveLive/dao/rdb"
	"LiveLive/model"
	"LiveLive/viper"
	"fmt"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

var Cfg *Config

var DB *gorm.DB
var Rdb *redis.Client

func Init() {

	con := viper.Init("db")

	var config Config
	if err := con.Viper.Unmarshal(&config); err != nil {
		log.Errorf("反序列化配置失败: %w", err)
	}

	Cfg = &config

	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port),
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.MySQL.User, Cfg.MySQL.Password, Cfg.MySQL.Host, Cfg.MySQL.Port, Cfg.MySQL.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	log.Info("数据库连接成功")

	model.MigrateUser(DB)
	model.MigrateLive(DB)
	model.MigrateCourseAndCourseMember(DB)
	model.MigrateCourseInvite(DB)
	model.MigrateQuestion(DB)

	db.Mysql = DB
	dao.Redis = Rdb
}
