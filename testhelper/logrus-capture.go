package testhelper

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

//LogCapturer logger interface
type LogCapturer interface {
	Release()
}

type logCapturer struct {
	*testing.T
	origOut       io.Writer
	jsonFormatter *logrus.TextFormatter
}

func (tl logCapturer) Format(entry *logrus.Entry) (rslt []byte, err error) {
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.PanicLevel {
		tl.Error(entry.Message)
	} else if entry.Level == logrus.FatalLevel {
		tl.Fatal(entry.Message)
	} else {
		tl.Logf(entry.Message)
	}
	return tl.jsonFormatter.Format(entry)
}
func (tl logCapturer) Write(p []byte) (n int, err error) {
	tl.Logf((string)(p))
	return len(p), nil
}

func (tl logCapturer) Release() {
	logrus.SetOutput(tl.origOut)
}

// CaptureLog redirects logrus output to testing.Log
func CaptureLog(t *testing.T) LogCapturer {
	fmt := logrus.TextFormatter{}
	lc := logCapturer{T: t, origOut: logrus.StandardLogger().Out, jsonFormatter: &fmt}
	logrus.SetOutput(lc)

	if !testing.Verbose() {
		logrus.SetOutput(lc)
	}
	return &lc
}
