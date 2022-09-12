package settings

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type KLogger struct {
	Name       string
	infologger *zap.SugaredLogger
	errlogger  *zap.SugaredLogger
}

func NewKLogger(name string) *KLogger {
	var err error
	id := strings.ReplaceAll(name, ":", "_")
	logpath := instance.LogPath
	var outpath string = ""

	if logpath != "none" && logpath != "" {
		outpath = Join(logpath, id)
		_, err = os.Stat(outpath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(outpath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating log path: %v", err)
			}
		}
	}

	var infowriter lumberjack.Logger
	var errwriter lumberjack.Logger

	if outpath != "" {
		infowriter = lumberjack.Logger{
			Filename:   Join(outpath, "info.log"),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
			Compress:   true,
		}

		errwriter = lumberjack.Logger{
			Filename:   Join(outpath, "error.log"),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
			Compress:   true,
		}
	}

	return &KLogger{Name: name, infologger: logInit(&infowriter), errlogger: logInit(&errwriter)}
}

func (l *KLogger) Info(msg string, v ...interface{}) {
	l.infologger.Infof("%s %v", msg, v)
}
func (l *KLogger) Error(msg string, v ...interface{}) {
	l.errlogger.Errorf("%s %v", msg, v)
}

func logInit(f *lumberjack.Logger) *zap.SugaredLogger {

	level := instance.GetZapLevel()

	pe := zap.NewProductionEncoderConfig()

	pe.EncodeTime = zapcore.ISO8601TimeEncoder // The encoder can be customized for each output
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	if f == nil {
		consolecore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
		return zap.New(consolecore).Sugar()
	}

	fileEncoder := zapcore.NewJSONEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	return zap.New(core).Sugar()

}
