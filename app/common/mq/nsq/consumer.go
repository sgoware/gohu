package nsq

import (
	"github.com/nsqio/go-nsq"
	"main/app/common/log"
	"time"
)

func NewConsumer(topic string, channel string) (consumer *nsq.Consumer, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	config.LookupdPollInterval = 15 * time.Second
	consumer, err = nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	zapLogger := log.GetSugaredLogger()
	logger := NewLogger(zapLogger)
	for i := 0; i <= nsq.LogLevelMax; i++ {
		consumer.SetLogger(logger, nsq.LogLevel(i))
	}

	h := &PrintHandler{Title: "print"}
	consumer.AddHandler(h)
	err = consumer.ConnectToNSQLookupds(MustGetNSQLookupAddrs())
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
