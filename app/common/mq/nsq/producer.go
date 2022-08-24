package nsq

import (
	"github.com/nsqio/go-nsq"
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
	err = producer.Ping()
	if err != nil {
		return nil, err
	}
	return
}
