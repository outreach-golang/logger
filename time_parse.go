package logger

import "time"

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
// 2020-11-08T08:18:46+08:00 => 2020-11-08 08:18:46
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	return ts.In(cst).Format(CSTLayout), nil
}

// CSTLayoutString 格式化时间
// 返回 "2006-01-02 15:04:05" 格式的时间
func CSTLayoutString() string {
	ts := time.Now()
	return ts.In(cst).Format(CSTLayout)
}

// CSTLayoutStringToUnix 返回 unix 时间戳
// 2020-01-24 21:11:11 => 1579871471
func CSTLayoutStringToUnix(cstLayoutString string) (int64, error) {
	stamp, err := time.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}
