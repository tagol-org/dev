/**
* @Author: tc.lam
* @Date: 2022/7/10 下午11:33
* @Software : GoLand
* @File: logger
* @Description:
* @Return:
**/

package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"log/syslog"
	"sync"
)

const AllLoggers = ""

type Logger struct {
	*logrus.Logger
	name    string
	lock    sync.Mutex
	loggers map[string]*Logger

	// keep output copy for closing it
	hooks []Hook
}

func (l *Logger) SetOutput(name string, output io.Writer) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if name == AllLoggers {
		l.Logger.Out = output
		for _, logger := range l.loggers {
			logger.SetOutput(AllLoggers, output)
		}
	} else {
		l.getLogger(name).SetOutput(AllLoggers, output)
	}
	return l
}

func (l *Logger) SetFormatter(name string, formatter Formatter) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if name == AllLoggers {
		l.Logger.Formatter = formatter
		for _, logger := range l.loggers {
			logger.SetFormatter(AllLoggers, formatter)
		}
	} else {
		l.getLogger(name).SetFormatter(AllLoggers, formatter)
	}
	return l
}

func (l *Logger) newZChildLogger(name string) *Logger {
	logger := NewLogger(name)
	logger.SetOutput(AllLoggers, l.Out).SetFormatter(AllLoggers, l.Formatter)

	return logger
}

func (l *Logger) getLogger(name string) *Logger {
	if name == AllLoggers {
		return l
	}
	if logger, ok := l.loggers[name]; ok {
		return logger
	}
	logger := l.newChildLogger(name)
	l.loggers[name] = logger
	return logger
}

func NewLogger(name string) *Logger {
	return &Logger{
		Logger:  logrus.New(),
		name:    name,
		loggers: make(map[string]*Logger),
	}
}
