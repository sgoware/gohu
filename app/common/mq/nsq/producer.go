package nsq

import (
	"github.com/nsqio/go-nsq"
	"main/app/common/log"
)

func NewProducer() (producer *nsq.Producer, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	producer, err = nsq.NewProducer(MustGetNSQDAddr(), config)
	if err != nil {
		return nil, err
	}

	zapLogger := log.GetSugaredLogger()
	logger := NewLogger(zapLogger)
	for i := 0; i <= nsq.LogLevelMax; i++ {
		producer.SetLogger(logger, nsq.LogLevel(i))
	}
	producer.SetLoggerLevel(nsq.LogLevelInfo)

	err = producer.Ping()
	if err != nil {
		return nil, err
	}

	err = producer.Publish("test", []byte("sdf"))
	if err != nil {
		return nil, err
	}

	return
}
