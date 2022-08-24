package log

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type ZapWriter struct {
	logger *zap.Logger
}

func GetLogXConfig(serviceName, level string) logx.LogConf {
	return logx.LogConf{
		ServiceName: serviceName,
		Mode:        "console",
		Encoding:    "plain",
		TimeFormat:  "2006-01-02 15:04:05.000",
		Level:       level,
	}
}

func GetZapWriter() (logx.Writer, error) {
	if logger == nil {
		return nil, errors.New("logger is null, try user NewLogger to initialize a logger")
	}
	return &ZapWriter{
		logger: logger,
	}, nil
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
