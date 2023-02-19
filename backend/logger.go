package backend

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var WireLogger *zap.Logger
var BackendLogger *zap.Logger
var FrontendLogger *zap.Logger

var encoderConfig zapcore.EncoderConfig = zapcore.EncoderConfig{
	MessageKey:   "message",
	CallerKey:    "from",
	EncodeCaller: zapcore.ShortCallerEncoder,
	TimeKey:      "when",
	EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC1123),
}

func newLogger(name string) *zap.Logger {
	// eventually the logs should live in a directory the app creates and manages, with rotation of some sort, but that is not a now problem
	config := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{fmt.Sprintf("%s.log", name)},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

// create 3 loggers
// one specifically for what comes off the websocket wire socket.log
// one for backend stuff go.log
// one for frontend stuff svelte.log
func init() {
	WireLogger = newLogger("socket")
	BackendLogger = newLogger("go")
	FrontendLogger = newLogger("svelte")
}
