package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Msg *zap.SugaredLogger

func init() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(consoleEncoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(log.Writer())), zapcore.DebugLevel)

	logger := zap.New(core)

	Msg = logger.Sugar()

	defer Msg.Sync()
}
