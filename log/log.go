package log

import (
	"os"
	"path/filepath"
	"seed-sync/config"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log   *zap.Logger        // 暴露原始 Logger
	Sugar *zap.SugaredLogger // 暴露 SugaredLogger
	once  sync.Once
)

// InitLogger 初始化日志
func InitLogger() {
	once.Do(func() {
		// 创建自定义 encoder
		encoder := getEncoder()

		// 创建 writeSyncer
		writeSyncer := getWriteSyncer()

		// 创建 core
		core := zapcore.NewCore(
			encoder,
			writeSyncer,
			zapcore.InfoLevel,
		)

		// 创建 logger
		Log = zap.New(
			core,
			zap.AddCaller(),      // 添加调用者信息
			zap.AddCallerSkip(1), // 跳过 1 层调用者
		)

		// 创建 sugar logger
		Sugar = Log.Sugar()

		// 替换全局 logger
		zap.ReplaceGlobals(Log)
	})
}

// getEncoder 自定义日志格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder, // 使用自定义的级别编码器
		EncodeTime:     customTimeEncoder,  // 自定义时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// customLevelEncoder 自定义级别编码器
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getWriteSyncer 获取日志写入器
func getWriteSyncer() zapcore.WriteSyncer {
	logConfig := config.Conf.LogConfig

	// 确保日志目录存在
	if err := os.MkdirAll(logConfig.Path, 0755); err != nil {
		panic("can't create log directory, err: " + err.Error())
	}

	// 设置日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logConfig.Path, time.Now().Format("2006-01-02")+".log"),
		MaxSize:    logConfig.MaxSize,    // 每个日志文件最大尺寸（MB）
		MaxBackups: logConfig.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     logConfig.MaxAge,     // 保留旧文件的最大天数
		Compress:   logConfig.Compress,   // 是否压缩/归档旧文件
	}

	// 同时输出到控制台和文件
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberJackLogger),
	)
}

// Debug logger
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info logger
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn logger
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error logger
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal logger
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
