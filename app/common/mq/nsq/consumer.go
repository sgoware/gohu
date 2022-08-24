package nsq

import (
	"github.com/nsqio/go-nsq"
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

	h := &PrintHandler{Title: "print"}
	consumer.AddHandler(h)

	err = consumer.ConnectToNSQLookupds(MustGetNSQLookupAddrs())
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
