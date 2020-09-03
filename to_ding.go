package logger

import (
	"bytes"
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
)

func WriteDing(l zapcore.Level, z zapcore.Encoder, configs *Config) zapcore.Core {

	r := &toDing{configs: configs}

	zlvl := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l
	})

	return zapcore.NewCore(z, zapcore.AddSync(r), zlvl)
}

type toDing struct {
	configs *Config
	TraceId string `json:"traceid"`
	Msg     string `json:"msg"`
}

func (t *toDing) Write(p []byte) (n int, err error) {

	msgToMap(t, p)

	var (
		errorMsg         = msgSplit(t.Msg)
		tid              = t.TraceId
		data             = make(map[string]string)
		currentTime      = time.Now().Format("2006-01-02 15:04:05")
		sendDataTemplete = `
#### 项目名称: ` + t.configs.ServerName + `
> 错误信息: ` + errorMsg + "\n" + `
> 机器地址: ` + LocalIP() + "\n" + `
> TraceId：[` + tid + "](https://sls.console.aliyun.com/lognext/project/test-prject-02/logsearch/test-log?traceid=" + tid + ")\n" + `
> 时间 : ` + currentTime
	)

	data["webhook"] = t.configs.Ding.DingWebhook

	data["content"] = sendDataTemplete
	data["title"] = "错误信息"

	jsonData, _ := jsoniter.Marshal(&data)

	client := &http.Client{}

	cxt, _ := context.WithTimeout(context.Background(), time.Second*2)

	req, err := http.NewRequest("POST", t.configs.Ding.DingHost, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req = req.WithContext(cxt)

	_, err = client.Do(req)

	if err != nil {
		return 0, err
	}
	return 0, nil
}

func msgToMap(t *toDing, p []byte) {
	_ = jsoniter.Unmarshal(p, t)
	return
}

func msgSplit(msg string) (rb string) {

	if len(msg) >= 500 {

		rb = msg[:500]

		return
	}

	return msg

}
