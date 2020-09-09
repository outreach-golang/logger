package logger

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
)

type GormLogger struct {
}

func (logger *GormLogger) Print(values ...interface{}) {

	if values[0] == "sql" {
		ctx := context.Background()
		NewContext(
			&ctx,
			zap.String("sql", values[3].(string)),
			zap.String("params", fmt.Sprint(values[4])),
			zap.String("rows.affected", fmt.Sprint(values[5])),
			zap.String("file.with.line.num", values[1].(string)),
			zap.String("sql.duration", fmt.Sprint(values[2].(time.Duration))),
		)

		if values[2].(time.Duration) >= 2000 {
			WithContext(&ctx).Error("sql：**" + values[3].(string) + "** 过慢")
		}

	}

	var (
		ctx     = context.Background()
		traceId = NewTraceID()
	)

	switch values[0] {
	case "sql":
		NewContext(
			&ctx,
			zap.String("traceid", traceId),
			zap.String("sql", values[3].(string)),
			zap.String("params", fmt.Sprint(values[4])),
			zap.String("rows.affected", fmt.Sprint(values[5])),
			zap.String("file.with.line.num", values[1].(string)),
			zap.String("sql.duration", fmt.Sprint(values[2].(time.Duration))),
		)

		if values[2].(time.Duration) >= (time.Millisecond * 2000) {
			WithContext(&ctx).Error("**sql：" + values[3].(string) + "\n参数：" +
				fmt.Sprint(values[4]) + "\n耗时：[" + fmt.Sprint(values[2].(time.Duration)) + "]** 过慢")
		}

		break
	case "log":
		NewContext(
			&ctx,
			zap.String("traceid", traceId),
			zap.String("sql.error", values[2].(*mysql.MySQLError).Message),
			zap.String("file.with.line.num", values[1].(string)),
		)
		WithContext(&ctx).Error("**" + values[2].(*mysql.MySQLError).Message + "**")
		break
	}

}
