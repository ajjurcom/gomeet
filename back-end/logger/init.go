package logger

import (
	"com/mittacy/gomeet/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

// InitLogger 初始化日志器，返回类型 error，如果返回错误，必须 panic
func InitLogger() error {
	writer, err := getLogWriter()
	if err != nil {
		return err
	}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	return nil
}

// Record 记录日志函数
func Record(msg string, errList ...error) {
	if len(errList) == 0 {
		sugarLogger.Infof("%s", msg)
	} else {
		sugarLogger.Errorf("%s\nerr: %s", msg, errList)
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() (zapcore.WriteSyncer, error) {
	// 1. 配置信息
	var (
		path string
		maxSize int
		maxBackups int
		maxAge int
		isCompress bool
		err error
	)
	// 2. 从配置文件获取信息
	logConfig := config.Cfg.Section("logs")
	path = logConfig.Key("record_log").String()
	if maxSize, err = logConfig.Key("max_size").Int(); err != nil {
		return nil, err
	}
	if maxBackups, err = logConfig.Key("max_backups").Int(); err != nil {
		return nil, err
	}
	if maxAge, err = logConfig.Key("max_age").Int(); err != nil {
		return nil, err
	}
	if isCompress, err = logConfig.Key("is_compress").Bool(); err != nil {
		return nil, err
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename: path,
		MaxSize: maxSize,	// 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackups,	// 保留旧文件的最大个数
		MaxAge: maxAge,		// 保留旧文件的最大天数
		Compress: isCompress,	// 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger), nil
}
