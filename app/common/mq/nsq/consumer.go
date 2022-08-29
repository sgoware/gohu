package nsq

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"main/app/common/log"
	"time"
)

type ConsumerService struct {
	consumer *nsq.Consumer
	logger   *zap.SugaredLogger
}

func (m *ConsumerService) Start() {
	err := m.consumer.ConnectToNSQLookupds(MustGetNSQLookupAddrs())
	if err != nil {
		m.logger.Errorf("start nsq consumer service failed, err: %v", err)
	}
}

func (m *ConsumerService) Stop() {
	m.Stop()
}

func NewConsumerService(topic string, channel string, handler nsq.Handler) (service *ConsumerService, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	config.LookupdPollInterval = 15 * time.Second
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	zapLogger := log.GetSugaredLogger()
	logger := NewLogger(zapLogger)
	for i := 0; i <= nsq.LogLevelMax; i++ {
		consumer.SetLogger(logger, nsq.LogLevel(i))
	}

	consumer.AddHandler(handler)

	return &ConsumerService{consumer: consumer, logger: zapLogger}, nil
}
