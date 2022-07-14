package log

import (
	"bytes"
	"testing"

	"github.com/tagol-org/dev/errors"
)

func TestLogger_SetOutput(t *testing.T) {
	l := NewLogger("test")
	output := bytes.NewBuffer(nil)
	logger1 := l.GetLogger("1")

	// set output for all loggers
	l.SetOutput("", output)
	logger1.Error("test")
	if output.Len() == 0 {
		t.Fatal("should log")
	}
	output.Reset()

	// set output for existing logger
	output1 := bytes.NewBuffer(nil)
	l.SetOutput("1", output1)
	logger1.Error("test")
	if output.Len() > 0 || output1.Len() == 0 {
		t.Fatal("output err")
	}
	output1.Reset()

	// set output for not existing logger
	output2 := bytes.NewBuffer(nil)
	l.SetOutput("2", output2)
	logger2 := l.GetLogger("2")
	logger2.Error("test")
	if output.Len() > 0 || output1.Len() > 0 || output2.Len() == 0 {
		t.Fatal("output err")
	}
}

func makeError() error {
	return errors.New("makeError")
}

func TestLoggerErrorStack(t *testing.T) {
	output := bytes.NewBuffer(nil)
	l := NewLoggerWithHooks("test", NewErrorStackHook(true)).SetFormatter(AllLoggers, ColoredTextFormatter)
	l.SetOutput("", output)
	err := makeError()
	l.WithError(errors.Annotate(err, "annotate")).Error("testerrorstack")
	if output.Len() == 0 {
		t.Fatal("should log")
	}
}
