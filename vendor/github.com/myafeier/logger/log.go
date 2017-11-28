package logger

import (
	"fmt"
	"io"
	"log"
)

type LogLevel int

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_WARNING
	LOG_ERR
	LOG_OFF
	LOG_UNKNOWN
)

const (
	DEFAULT_LOG_PREFIX = "[wss]"
	DEFAULT_LOG_FLAG   = log.Lshortfile | log.Ldate | log.Lmicroseconds
	DEFAULT_LOG_LEVEL  = LOG_DEBUG
)

type ILogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Level() LogLevel
	SetLevel(l LogLevel)
}

// SimpleLogger is the default implement of ILogger
type SimpleLogger struct {
	DEBUG   *log.Logger
	ERR     *log.Logger
	INFO    *log.Logger
	WARN    *log.Logger
	level   LogLevel
	showSQL bool
}

var _ ILogger = &SimpleLogger{}

// NewSimpleLogger use a special io.Writer as logger output
func NewSimpleLogger(out io.Writer) *SimpleLogger {
	return NewSimpleLogger2(out, DEFAULT_LOG_PREFIX, DEFAULT_LOG_FLAG)
}

// NewSimpleLogger2 let you customrize your logger prefix and flag
func NewSimpleLogger2(out io.Writer, prefix string, flag int) *SimpleLogger {
	return NewSimpleLogger3(out, prefix, flag, DEFAULT_LOG_LEVEL)
}

// NewSimpleLogger3 let you customrize your logger prefix and flag and logLevel
func NewSimpleLogger3(out io.Writer, prefix string, flag int, l LogLevel) *SimpleLogger {
	return &SimpleLogger{
		DEBUG: log.New(out, fmt.Sprintf("%s %c[%d;%d;%dm [debug] %c[0m ", prefix, 0x1B, 1, 40, 32, 0x1B), flag),
		ERR:   log.New(out, fmt.Sprintf("%s %c[%d;%d;%dm [error] %c[0m ", prefix, 0x1B, 1, 40, 31, 0x1B), flag),
		INFO:  log.New(out, fmt.Sprintf("%s %c[%d;%d;%dm [info] %c[0m ", prefix, 0x1B, 1, 40, 34, 0x1B), flag),
		WARN:  log.New(out, fmt.Sprintf("%s %c[%d;%d;%dm [warn] %c[0m ", prefix, 0x1B, 1, 40, 33, 0x1B), flag),
		level: l,
	}
}

// Error implement core.ILogger
func (s *SimpleLogger) Error(v ...interface{}) {
	if s.level <= LOG_ERR {
		s.ERR.Output(2, fmt.Sprint(v...))
	}
	return
}

// Errorf implement core.ILogger
func (s *SimpleLogger) Errorf(format string, v ...interface{}) {
	if s.level <= LOG_ERR {
		s.ERR.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Debug implement core.ILogger
func (s *SimpleLogger) Debug(v ...interface{}) {
	if s.level <= LOG_DEBUG {
		s.DEBUG.Output(2, fmt.Sprint(v...))
	}
	return
}

// Debugf implement core.ILogger
func (s *SimpleLogger) Debugf(format string, v ...interface{}) {
	if s.level <= LOG_DEBUG {
		s.DEBUG.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Info implement core.ILogger
func (s *SimpleLogger) Info(v ...interface{}) {
	if s.level <= LOG_INFO {
		s.INFO.Output(2, fmt.Sprint(v...))
	}
	return
}

// Infof implement core.ILogger
func (s *SimpleLogger) Infof(format string, v ...interface{}) {
	if s.level <= LOG_INFO {
		s.INFO.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Warn implement core.ILogger
func (s *SimpleLogger) Warn(v ...interface{}) {
	if s.level <= LOG_WARNING {
		s.WARN.Output(2, fmt.Sprint(v...))
	}
	return
}

// Warnf implement core.ILogger
func (s *SimpleLogger) Warnf(format string, v ...interface{}) {
	if s.level <= LOG_WARNING {
		s.WARN.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Level implement core.ILogger
func (s *SimpleLogger) Level() LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *SimpleLogger) SetLevel(l LogLevel) {
	s.level = l
	return
}
