package logger

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	gl "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Writer interface {
	Printf(string, ...interface{})
}

type GormLoggerV2 struct {
	Writer
	SlowSqlTime time.Duration
	LogLevel    gl.LogLevel
}

func (l *GormLoggerV2) LogMode(level gl.LogLevel) gl.Interface {
	newLogger := &GormLoggerV2{}
	newLogger.LogLevel = level
	return newLogger
}

func (l *GormLoggerV2) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= gl.Info {
		fmt.Print(ctx)
		l.Printf(s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (l *GormLoggerV2) Warn(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (l *GormLoggerV2) Error(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (l *GormLoggerV2) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	panic("implement me")
}

func (l *GormLoggerV2) Print(values ...interface{}) {

	var (
		ctx     = context.Background()
		traceId = NewTraceID()
	)

	switch values[0] {
	case "sql":
		newContext := NewContext(
			ctx,
			zap.String("traceid", traceId),
			zap.String("sql", values[3].(string)),
			zap.String("params", fmt.Sprint(values[4])),
			zap.String("rows.affected", fmt.Sprint(values[5])),
			zap.String("file.with.line.num", values[1].(string)),
			zap.String("sql.duration", fmt.Sprint(values[2].(time.Duration))),
		)

		if l.SlowSqlTime > 0 && values[2].(time.Duration) >= (time.Millisecond*l.SlowSqlTime) {
			WithContext(newContext).Error("**sql：" + values[3].(string) + "\n参数：" +
				fmt.Sprint(values[4]) + "\n耗时：[" + fmt.Sprint(values[2].(time.Duration)) + "]** 过慢")
		} else {
			WithContext(newContext).Info("sql：" + values[3].(string) + "\n参数：" +
				fmt.Sprint(values[4]) + "\n耗时：[" + fmt.Sprint(values[2].(time.Duration)) + "]")
		}

		break
	case "log":

		var sqlErrorField string
		switch values[2].(type) {
		case *mysql.MySQLError:
			sqlErrorField = values[2].(*mysql.MySQLError).Message
			break
		default:
			sqlErrorField = values[2].(error).Error()
		}

		newContext := NewContext(
			ctx,
			zap.String("traceid", traceId),
			zap.String("sql.error", sqlErrorField),
			zap.String("file.with.line.num", values[1].(string)),
		)
		WithContext(newContext).Error("**" + sqlErrorField + "**")
		break
	}

}
