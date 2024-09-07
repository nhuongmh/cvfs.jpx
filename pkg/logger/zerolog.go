package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {

}

func InitLog() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	multi := zerolog.MultiLevelWriter(consoleWriter)
	Log = zerolog.New(multi).With().Timestamp().Logger()
}

func InitLogWithWriter(logCfg ...io.Writer) {
	if logCfg == nil {
		InitLog()
		return
	}

	multi := zerolog.MultiLevelWriter(logCfg...)
	Log = zerolog.New(multi).With().Timestamp().Logger()
}
