package log

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tagol-org/dev/errors"
)

const LoggerNameField = "logger"

type LoggerNameHook struct {
	tag string
}

func (*LoggerNameHook) Levels() []Level {
	return AllLevels
}

func (h *LoggerNameHook) Fire(entry *Entry) error {
	entry.Data[LoggerNameField] = h.tag
	return nil
}

func NewLoggerNameHook(tag string) *LoggerNameHook {
	return &LoggerNameHook{tag: tag}
}

type RuntimeHook struct{}

func (h *RuntimeHook) Levels() []Level {
	return AllLevels
}

func (h *RuntimeHook) Fire(entry *Entry) error {
	file := "???"
	funcName := "???"
	line := 0

	pc := make([]uintptr, 64)
	// Skip runtime.Callers, self, and another call from logrus
	n := runtime.Callers(3, pc)
	if n != 0 {
		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		// Loop to get frames.
		// A fixed number of pcs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()
			if !strings.Contains(frame.File, "github.com/sirupsen/logrus") && !strings.Contains(frame.Function, "git.in.chaitin.net/dev/go") {
				file = frame.File
				funcName = frame.Function
				line = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}

	slices := strings.Split(file, "/")
	file = slices[len(slices)-1]

	entry.Data["file"] = file
	entry.Data["func"] = funcName
	entry.Data["line"] = line
	return nil
}

func NewRuntimeHook() *RuntimeHook {
	return &RuntimeHook{}
}

type ErrorStackHook struct {
	appendToMessage bool
}

func (*ErrorStackHook) Levels() []Level {
	return AllLevels
}

func (e *ErrorStackHook) Fire(entry *Entry) error {
	if err, ok := entry.Data[logrus.ErrorKey].(error); ok {
		errstack := errors.ErrorStack(err)
		if e.appendToMessage {
			entry.Message += "\n"
			entry.Message += errstack
			entry.Message += "\n"
		} else {
			entry.Data["errorstack"] = errstack
		}
	}
	return nil
}

func NewErrorStackHook(appendToMessage bool) *ErrorStackHook {
	return &ErrorStackHook{appendToMessage: appendToMessage}
}
