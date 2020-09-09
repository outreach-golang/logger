package logger

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

type SaveLogForm int

const (
	GinKey                = "gin_logger"
	NormalKey             = "normal_logger"
	File      SaveLogForm = iota
	Ding
	AliLog
)

var GLogger *zap.Logger

func NewLogger(ops ...Option) (*zap.Logger, error) {

	config := DefaultConfig()

	for _, op := range ops {
		op(config)
	}

	if config.ServerName == "" {
		return nil, errors.New("ServerName参数必填！")
	}

	Logger, err := logger(config)
	if err != nil {
		return nil, err
	}

	GLogger = Logger

	return Logger, err
}

func logger(configs *Config) (*zap.Logger, error) {

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= configs.EnableLogLevel
	})

	var core zapcore.Core

	switch configs.SaveLogAddr {
	case File:
		infoWriter := GetWriter(configs)
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		)
	case Ding:
		core = zapcore.NewTee(
			WriteDing(zap.ErrorLevel, encoder, configs),
		)
	case AliLog:
		core = zapcore.NewTee(
			WriteAliLog(configs.EnableLogLevel, encoder, configs),
			WriteDing(zap.ErrorLevel, encoder, configs),
		)
	}

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(infoLevel))

	return log, nil
}

func GetWriter(configs *Config) io.Writer {

	hook, err := rotatelogs.New(
		strings.Replace(configs.File.FilePath, ".log", "", -1)+"-%Y%m%d.txt",
		rotatelogs.WithMaxAge(configs.File.FileMaxAge),
		rotatelogs.WithRotationTime(configs.File.FileRotationTime),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func NewContext(ctx context.Context, fields ...zapcore.Field) {
	if ctx == nil {
		ctx = context.Background()
	}

	if gc, ok := ctx.(*gin.Context); ok {
		gc.Set(GinKey, WithContext(ctx).With(fields...))
	} else {
		ctx = context.WithValue(ctx, NormalKey, WithContext(ctx).With(fields...))
	}

}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return GLogger
	}

	if gc, ok := ctx.(*gin.Context); ok {
		l, _ := gc.Get(GinKey)

		logger, ok := l.(*zap.Logger)

		if ok {
			return logger
		}

	} else {
		logger, ok := ctx.Value(NormalKey).(*zap.Logger)

		if ok {
			return logger
		}
	}

	return GLogger
}
