package logger

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func WriteAliLog(l zapcore.Level, z zapcore.Encoder, configs *Config) zapcore.Core {

	r := &toAliLog{configs: configs}
	r.initServ()

	zlvl := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l
	})

	return zapcore.NewCore(z, zapcore.AddSync(r), zlvl)
}

type toAliLog struct {
	configs *Config
	Client  *producer.Producer
}

func (t *toAliLog) initServ() {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = t.configs.AliLog.Endpoint
	producerConfig.AccessKeyID = t.configs.AliLog.AccessKeyID
	producerConfig.AccessKeySecret = t.configs.AliLog.AccessKeySecret
	t.Client = producer.InitProducer(producerConfig)
	t.Client.Start()
}

func (t *toAliLog) Write(p []byte) (n int, err error) {

	var (
		log     = &sls.Log{}
		content []*sls.LogContent
		data    = make(map[string]string)
	)

	if err = jsoniter.Unmarshal(p, &data); err != nil {
		return 0, err
	}

	for k, v := range data {
		content = append(content, &sls.LogContent{
			Key:   proto.String(k),
			Value: proto.String(v),
		})
	}

	log.Time = proto.Uint32(uint32(time.Now().Unix()))
	log.Contents = content

	err = t.Client.SendLog(t.configs.AliLog.Project, t.configs.AliLog.LogStore, t.configs.AliLog.Topic, t.configs.AliLog.Source, log)
	fmt.Println("发送日志报错：", err)

	return 0, nil
}
