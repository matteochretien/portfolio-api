package zap

import (
	"github.com/Niromash/niromash-api/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() (logger.Logger, error) {
	config := zap.NewDevelopmentConfig()
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeCaller = nil
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02/01/2006 - 15:04:05")
	config.EncoderConfig = encoderConfig

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger.Sugar(), nil
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}

func (z *ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.logger.Debugw(msg, keysAndValues...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

func (z *ZapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.logger.Fatalw(msg, keysAndValues...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

func (z *ZapLogger) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}

func (z *ZapLogger) Panicf(template string, args ...interface{}) {
	z.logger.Panicf(template, args...)
}

func (z *ZapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.logger.Panicw(msg, keysAndValues...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z *ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}
