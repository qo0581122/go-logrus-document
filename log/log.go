package Log

import (
	"os"

	. "./formatter"

	"github.com/sirupsen/logrus"
)

var Logger = NewLog()

type Log struct {
	log *logrus.Logger
}

func NewLog() *Log {
	mLog := logrus.New()
	mLog.SetOutput(os.Stderr)
	mLog.SetReportCaller(true)
	mLog.SetFormatter(&LogFormatter{})
	mLog.SetLevel(logrus.DebugLevel)
	return &Log{
		log: mLog,
	}
}

func (l *Log) Debug(args ...interface{}) {
	l.log.Debugln(args...)
}
func (l *Log) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}
func (l *Log) Info(args ...interface{}) {
	l.log.Infoln(args...)
}
func (l *Log) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}
func (l *Log) Error(args ...interface{}) {
	l.log.Errorln(args...)
}
func (l *Log) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
func (l *Log) Trace(args ...interface{}) {
	l.log.Traceln()
}
func (l *Log) Tracef(format string, args ...interface{}) {
	l.log.Tracef(format, args...)
}
func (l *Log) Panic(args ...interface{}) {
	l.log.Panicln()
}
func (l *Log) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *Log) Print(args ...interface{}) {
	l.log.Println()
}
func (l *Log) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}
