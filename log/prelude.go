package log

import "github.com/sirupsen/logrus"

const (
	PanicLevel Level = logrus.PanicLevel
	FatalLevel Level = logrus.FatalLevel
	ErrorLevel Level = logrus.ErrorLevel
	WarnLevel  Level = logrus.WarnLevel
	InfoLevel  Level = logrus.InfoLevel
	DebugLevel Level = logrus.DebugLevel
	TraceLevel Level = logrus.TraceLevel
)

type (
	Level         = logrus.Level
	Entry         = logrus.Entry
	Formatter     = logrus.Formatter
	Fields        = logrus.Fields
	Hook          = logrus.Hook
	TextFormatter = logrus.TextFormatter
	JSONFormatter = logrus.JSONFormatter
)

var (
	AllLevels = logrus.AllLevels

	ColoredTextFormatter Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}
)
