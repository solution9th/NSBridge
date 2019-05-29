package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Log log
// 直接使用该方法会导致 caller 级别顺序不对
// utils.Log.Info().Str("aaa", "bbb").Msg("ookk")
var Log zerolog.Logger

// NewLogFile 初始化文件日志
func NewLogFile(filename string, timeFormat string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("Open log file %s failed: %s", filename, err.Error()))
	}
	NewLog(f, timeFormat)
}

// NewLog new log
func NewLog(w io.Writer, timeFormat string) {

	TimeFormatUnixNano := "2006-01-02 15:04:05.999999999"

	zerolog.TimeFieldFormat = TimeFormatUnixNano
	output := zerolog.ConsoleWriter{Out: w, TimeFormat: TimeFormatUnixNano, NoColor: false}

	Log = zerolog.New(output).With().Caller().Timestamp().CallerWithSkipFrameCount(3).Logger()

}

// Debugf debug format
func Debugf(format string, v ...interface{}) {

	Log.Debug().Msgf(format, v...)
}

// Infof info format
func Infof(format string, v ...interface{}) {

	Log.Info().Msgf(format, v...)

}

// Warnf warn format
func Warnf(format string, v ...interface{}) {

	Log.Warn().Msgf(format, v...)
}

// Errorf error format
func Errorf(format string, v ...interface{}) {

	Log.Error().Msgf(format, v...)
}

// Error error format
func Error(v ...interface{}) {

	s := ""
	for range v {
		s = s + "%v "
	}

	Log.Error().Msgf(s, v...)
}

// Fatalf fatal format
func Fatalf(format string, v ...interface{}) {

	Log.Fatal().Msgf(format, v...)
}
