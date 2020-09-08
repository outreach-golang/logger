package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type GormLogger struct {
}

func (logger *GormLogger) Print(values ...interface{}) {

	if values[0] == "sql" {
		ctx := context.Background()
		NewContext(
			ctx,
			zap.String("sql", values[3].(string)),
			zap.String("params", fmt.Sprint(values[4])),
			zap.String("rows.affected", fmt.Sprint(values[5])),
			zap.String("file.with.line.num", values[1].(string)),
			zap.String("sql.duration", values[2].(string)),
		)

		if values[2].(time.Duration) >= 2000 {
			WithContext(ctx).Error("sql：**" + values[3].(string) + "** 过慢")
		}

	}

}
