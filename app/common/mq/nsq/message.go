package nsq

import (
	"github.com/nsqio/go-nsq"
	"main/app/common/log"
)

type PrintHandler struct {
	Title string
}

func (h *PrintHandler) HandleMessage(msg *nsq.Message) (err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("%s recv from %v, msg:%v\n", h.Title, msg.NSQDAddress, string(msg.Body))
	return
}
