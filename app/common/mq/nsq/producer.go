package nsq

import (
	"github.com/nsqio/go-nsq"
	"main/app/common/log"
)

var producer *nsq.Producer

func NewProducer() (*nsq.Producer, error) {
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

	return producer, nil
}

func GetProducer() (*nsq.Producer, error) {
	if producer == nil {
		return NewProducer()
	}
	return producer, nil
}
