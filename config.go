package logger

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type Option func(*Config)

type Config struct {
	//一般为项目名字
	ServerName string
	//日志存储的地方
	SaveLogAddr SaveLogForm
	//日志最小存储等级
	EnableLogLevel zapcore.Level
	//文件配置
	File FileConfig
	//Ding配置
	Ding DingConfig
	//阿里log配置
	AliLog AliLogConfig
}

type AliLogConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Project         string
	LogStore        string
	Topic           string
	Source          string
}

type DingConfig struct {
	DingHost    string
	DingWebhook string
}

type FileConfig struct {
	FilePath         string
	FileMaxAge       time.Duration
	FileRotationTime time.Duration
}

func DefaultConfig() *Config {
	var c = Config{}

	c.SaveLogAddr = File
	c.EnableLogLevel = zapcore.InfoLevel
	c.File.FilePath = "./logs/log"
	c.File.FileMaxAge = time.Hour * 24 * 3
	c.File.FileRotationTime = time.Hour * 24
	c.AliLog.Source = LocalIP()
	c.AliLog.Topic = ""

	return &c
}

func Source(s string) Option {
	return func(config *Config) {
		config.AliLog.Source = s
	}
}

func Topic(s string) Option {
	return func(config *Config) {
		config.AliLog.Topic = s
	}
}

func LogStore(s string) Option {
	return func(config *Config) {
		config.AliLog.LogStore = s
	}
}

func Project(s string) Option {
	return func(config *Config) {
		config.AliLog.Project = s
	}
}

func AccessKeyID(s string) Option {
	return func(config *Config) {
		config.AliLog.AccessKeyID = s
	}
}

func AccessKeySecret(s string) Option {
	return func(config *Config) {
		config.AliLog.AccessKeySecret = s
	}
}

func Endpoint(s string) Option {
	return func(config *Config) {
		config.AliLog.Endpoint = s
	}
}

func SaveLogAddr(s SaveLogForm) Option {
	return func(config *Config) {
		config.SaveLogAddr = s
	}
}

func EnableLogLevel(e zapcore.Level) Option {
	return func(config *Config) {
		config.EnableLogLevel = e
	}
}

func FilePath(s string) Option {
	return func(config *Config) {
		config.File.FilePath = s
	}
}

func FileMaxAge(s time.Duration) Option {
	return func(config *Config) {
		config.File.FileMaxAge = s
	}
}

func FileRotationTime(s time.Duration) Option {
	return func(config *Config) {
		config.File.FileRotationTime = s
	}
}

func DingHost(s string) Option {
	return func(config *Config) {
		config.Ding.DingHost = s
	}
}

func DingWebhook(s string) Option {
	return func(config *Config) {
		config.Ding.DingWebhook = s
	}
}

func ServerName(s string) Option {
	return func(config *Config) {
		config.ServerName = s
	}
}
