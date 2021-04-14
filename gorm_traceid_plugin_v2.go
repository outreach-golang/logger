package logger

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/outreach-golang/logger/gorm_V2"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"time"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
	SlowSqlTime        = 0.5
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var (
	_ gorm.Plugin = &TracePlugin{}
)

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	_ctx := db.Statement.Context

	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	sqlInfo := &gorm_V2.SqlInfo{}
	sqlInfo.Set("Timestamp", CSTLayoutString())
	sqlInfo.Set("SQL", sql)
	sqlInfo.Set("Stack", utils.FileWithLineNum())
	sqlInfo.Set("Rows", db.Statement.RowsAffected)
	sqlInfo.Set("CostSeconds", time.Since(ts).Seconds())
	sqlInfo.Set("Table", time.Since(ts).Seconds())

	//sqlInfo := &gorm_V2.SQL{}
	//sqlInfo.Timestamp = CSTLayoutString()
	//sqlInfo.SQL = sql
	//sqlInfo.Stack = utils.FileWithLineNum()
	//sqlInfo.Rows = db.Statement.RowsAffected
	//sqlInfo.CostSeconds = time.Since(ts).Seconds()
	//sqlInfo.Table = db.Statement.Table

	switch db.Error {
	case nil:
	default:
		sqlJson, _ := jsoniter.MarshalToString(sqlInfo)

		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			WithContext(_ctx).Info(db.Error.Error(), zap.String("sql.info", sqlJson))
		} else {
			WithContext(_ctx).Error(db.Error.Error(), zap.Any("sql.info", sqlJson))
		}
	}

	if sqlInfo.Get("CostSeconds").(float64) >= SlowSqlTime {
		wbuff := buffer.Buffer{}
		wbuff.AppendString(sql)
		wbuff.AppendString("-----**执行时间：【 ")
		wbuff.AppendFloat(sqlInfo.Get("CostSeconds").(float64), 64)
		wbuff.AppendString(" 秒】**")

		WithContext(_ctx).Error(wbuff.String())
	}

	return
}
