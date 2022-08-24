package nsq

import (
	"errors"
	"go.uber.org/zap"
)

type logger struct {
	zap *zap.SugaredLogger
}

func (l *logger) Output(calldepth int, s string) error {
	if l == nil || l.zap == nil {
		return errors.New("empty logger")
	}
	l.zap.Debug(s)
	return nil
}

func NewLogger(zap *zap.SugaredLogger) *logger {
	return &logger{zap: zap}
}
