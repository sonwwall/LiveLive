package dao

import (
	"LiveLive/dao/db"
	"LiveLive/model"
	"LiveLive/viper"
	"fmt"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
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

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

var Cfg *Config

var DB *gorm.DB

func Init() {

	con := viper.Init("db")

	var config Config
	if err := con.Viper.Unmarshal(&config); err != nil {
		log.Errorf("反序列化配置失败: %w", err)
	}

	Cfg = &config

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

	db.Mysql = DB
}
