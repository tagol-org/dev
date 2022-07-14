package log

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
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

func (l *Logger) Name() string {
	return l.name
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

func (l *Logger) SetLevel(name string, level Level) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if name == AllLoggers {
		l.Logger.Level = Level(level)
		for _, logger := range l.loggers {
			logger.SetLevel(name, level)
		}
	} else {
		l.getLogger(name).SetLevel(AllLoggers, level)
	}
	return l
}

func (l *Logger) AddHook(name string, hook Hook) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if name == AllLoggers {
		l.hooks = append(l.hooks, hook)
		l.Logger.AddHook(hook)
		for _, logger := range l.loggers {
			logger.AddHook(AllLoggers, hook)
		}
	} else {
		l.getLogger(name).AddHook(AllLoggers, hook)
	}
	return l
}

func (l *Logger) newChildLogger(name string) *Logger {
	logger := NewLogger(name)
	logger.SetOutput(AllLoggers, l.Out).
		SetFormatter(AllLoggers, l.Formatter).
		SetLevel(AllLoggers, Level(l.Level))
	for _, hook := range l.hooks {
		logger.AddHook(AllLoggers, hook)
	}
	// LoggerNameHook only add to current logger
	logger.Logger.AddHook(NewLoggerNameHook(name))

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

func (l *Logger) GetLogger(name string) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.getLogger(name)
}

func NewLogger(name string) *Logger {
	return &Logger{
		Logger:  logrus.New(),
		name:    name,
		loggers: make(map[string]*Logger),
	}
}

func NewLoggerWithHooks(name string, hooks ...Hook) (logger *Logger) {
	logger = &Logger{
		Logger:  logrus.New(),
		name:    name,
		loggers: make(map[string]*Logger),
	}
	for _, hook := range hooks {
		logger.AddHook(AllLoggers, hook)
	}
	return
}
