package log

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const InfoLevel = "info"

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func setHandler(l zapcore.Level) {
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewDevelopmentConfig()

	if l.String() == InfoLevel {
		encoderCfg.DisableCaller = true
	}

	encoderCfg.EncoderConfig.TimeKey = "time"
	encoderCfg.EncoderConfig.EncodeTime = timeEncoder
	encoderCfg.EncoderConfig.LevelKey = "level"
	encoderCfg.EncoderConfig.MessageKey = "message"
	encoderCfg.EncoderConfig.CallerKey = "caller"
	encoderCfg.Level = atom
	encoderCfg.Encoding = "console"

	logger, err := encoderCfg.Build()
	if err != nil {
		log.Println(err)
	}

	zap.ReplaceGlobals(logger)
	atom.SetLevel(l)
}

func Verbosity(l string) {
	switch l {
	case "error":
		setHandler(zap.ErrorLevel)
	case "info":
		setHandler(zap.InfoLevel)
	case "debug":
		setHandler(zap.DebugLevel)
	default:
		setHandler(zap.DebugLevel)
	}
}
