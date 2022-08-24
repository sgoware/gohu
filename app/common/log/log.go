package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/app/utils"
	"os"
	"time"
)

var (
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

func NewLogger() *zap.Logger {
	options := Options{
		SavePath:     "log",
		EncoderType:  ConsoleEncode,
		EncodeLevel:  CapitalColorLevelEncoder,
		EncodeCaller: FullCallerEncoder,
	}
	return NewLoggerWithOptions(options)
}

func NewLoggerWithOptions(options Options) *zap.Logger {
	// 创建日志保存的文件夹
	if err := utils.IsNotExistMkDir(options.SavePath); err != nil {
		fmt.Printf("Create %v directory\n", options.SavePath)
		_ = os.Mkdir(options.SavePath, os.ModePerm)
	}
	dynamicLevel := zap.NewAtomicLevel()
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	encoder := getEncoder(options)
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, dynamicLevel), //控制台输出
		//日志文件输出,按等级归档
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/all/server_all.log", options.SavePath)), zapcore.DebugLevel),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/debug/server_debug.log", options.SavePath)), debugPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/info/server_info.log", options.SavePath)), infoPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/warn/server_warn.log", options.SavePath)), warnPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/error/server_error.log", options.SavePath)), errorPriority),
	}
	zapLogger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	defer zapLogger.Sync()
	// 将当前日志等级设置为Debug
	// 注意日志等级低于设置的等级，日志文件也不分记录
	dynamicLevel.SetLevel(zap.DebugLevel)
	//设置全局logger
	logger = zapLogger
	sugaredLogger = zapLogger.Sugar()
	logger.Info("Initialize logger successfully!")
	//sugar.Debug("test")
	//sugar.Warn("test")
	//sugar.Error("test")
	//sugar.DPanic("test")
	//sugar.Panic("test") //打印后程序停止,defer执行
	//sugar.Fatal("test") //打印后程序停止,defer不执行
	return logger
}

func GetLogger() *zap.Logger {
	if logger == nil {
		NewLogger()
	}
	return logger
}

func GetSugaredLogger() *zap.SugaredLogger {
	if sugaredLogger == nil {
		NewLogger()
	}
	return sugaredLogger
}

// 获取编码器
func getEncoder(options Options) zapcore.Encoder {
	if options.EncoderType == JsonEncoder {
		return zapcore.NewJSONEncoder(getEncoderConfig(options))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(options))
}

// 编码器设置
func getEncoderConfig(options Options) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",                     // 日志消息键
		LevelKey:       "level",                       // 日志等级键
		TimeKey:        "time",                        // 时间键
		NameKey:        "logger",                      // 日志记录器名
		CallerKey:      "caller",                      // 日志文件信息键
		StacktraceKey:  "stacktrace",                  // 堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 友好日志等级名大小写（info INFO）
		EncodeTime:     CustomTimeEncoder,             // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.FullCallerEncoder,     // 日志文件信息 short（包/文件.go:行号） full (文件位置.go:行号)
	}
	switch {
	case options.EncodeLevel == LowercaseLevelEncoder: // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case options.EncodeLevel == LowercaseColorLevelEncoder: // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case options.EncodeLevel == CapitalLevelEncoder: // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case options.EncodeLevel == CapitalColorLevelEncoder: // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	if options.EncodeCaller == ShortCallerEncoder {
		config.EncodeCaller = zapcore.ShortCallerEncoder
	}
	return config
}

// 读写器设置
func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    1,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 100,  // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// CustomTimeEncoder 格式化时间
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}
