package logger

const (

	// QISHUTimeFormat ...
	QISHUTimeFormat = "200601021504"

	// DefaultTimeFormat ...
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

// Logrotate ...
type Logrotate struct {

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" bson:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" bson:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" bson:"maxbackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" bson:"compress"`
}

// Config ...
type Config struct {

	// 日志输出类型，支持 DEFAULT、TEXT 和 JSON 三种。 默认为 DEFAULT
	Encoder string `json:"encoder" bson:"encoder"`

	// 日志级别，支持 DEBUG、INFO、WARN、ERROR 四种。 默认为 WARN
	Level string `json:"level" bson:"level"`

	// 是否开启函数名、文件行数打印
	ReportCaller bool `json:"report_caller" bson:"report_caller"`

	// 时间格式，支持 DEFAULT 和 QISHU 两种，默认为 DEFAULT
	TimeFormat string `json:"time_format" bson:"time_format"` // DEFAULT/QISHU

	// 日志输出，支持 syslog://HOST:PORT、STDOUT、STDERR 和 FILE 四种。默认为 STDOUT
	// 注: SYSLOG 仅支持 UDP
	Output []string `json:"output" bson:"output"`

	// SYSLOG 日志需要的 APP 字段， 默认为 go-kit
	App string `json:"app" bson:"app"`

	// FILE 日志需要的切分配置
	Logrotate *Logrotate `json:"logrotate" bson:"logrotate"`
}
