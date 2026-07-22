package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey struct{}

var key = loggerContextKey{}

type Logger struct {
	*zap.Logger
	file *os.File
}

func (logger *Logger) Close() {
	if logger == nil {
		return
	}

	_ = logger.Sync()

	if logger.file != nil {
		_ = logger.file.Close()
	}
}

func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("failed to get the Logger from the context")
	}
	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0o755); err != nil {
		return nil, fmt.Errorf("create folder: %w", err)
	}

	timeStamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timeStamp),
	)
	file, err := os.OpenFile(
		logFilePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(file), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   file,
	}, nil
}
