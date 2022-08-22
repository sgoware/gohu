package log

// Options 日志器设置
type Options struct {
	SavePath     string // 日志保存位置
	EncoderType  string // 日志编码器类型("json","console")
	EncodeLevel  string // 日志编码器风格
	EncodeCaller string // 打印调用函数风格
}

const (
	JsonEncoder   = "json"
	ConsoleEncode = "console"

	LowercaseLevelEncoder      = "LowercaseLevelEncoder"
	LowercaseColorLevelEncoder = "LowercaseColorLevelEncoder"
	CapitalLevelEncoder        = "CapitalLevelEncoder"
	CapitalColorLevelEncoder   = "CapitalColorLevelEncoder"

	ShortCallerEncoder = "ShortCallerEncoder"
	FullCallerEncoder  = "FullCallerEncoder"
)
