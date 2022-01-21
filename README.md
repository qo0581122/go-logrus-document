# go-logrus-document

### Go插件之logrus

> Logrus是Go (golang)的结构化日志程序，完全兼容标准库的API日志程序。
> Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.

**文章目录**:

- Logrus自带两种formatter
  - TextFormatter
  - JsonFormatter
  - 自定义Formatter
- Logrus基本用法
- 自定义Log
  - 包结构
  - formatter.go
  - log.go
  - main.go

**注意：基本用法请跳转[Logrus](https://github.com/sirupsen/logrus)**

---

### 1 Logrus自带两种formatter

#### 1.1 TextFormatter

```
下面展示几个常用的字段
type TextFormatter struct {
	DisableColors bool // 开启颜色显示
	
	DisableTimestamp bool // 开启时间显示

	TimestampFormat string	// 自定义时间格式

	QuoteEmptyFields bool	//空字段括在引号中

	CallerPrettyfier func(*runtime.Frame) (function string, file string) //用于自定义方法名和文件名的输出
}
```

#### 1.2  JsonFormatter

```
下面展示几个常用的字段
type JSONFormatter struct {
	TimestampFormat string // 自定义时间格式

	DisableTimestamp bool // 开启时间显示

	CallerPrettyfier func(*runtime.Frame) (function string, file string) //用于自定义方法名和文件名的输出

	PrettyPrint bool //将缩进所有json日志
}
```

#### 1.3 第三种 自定义Formatter

```go
只需要实现该接口
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

其中entry参数
type Entry struct {
	// Contains all the fields set by the user.
	Data Fields
	
	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	//Calling method, with package name
	Caller *runtime.Frame

	//Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic  
	Message string 
	
	//When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer
}
```

### 2 Logrus基本用法

```go
func Demo(log *logrus.Logger) {
	log.Info("i'm demo")

}

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04", //自定义日期格式
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) { //自定义Caller的返回
			//处理文件名
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})
	Demo(log)
}
```

### 3 自定义Log

#### 3.1 包结构

```
Test
  - log 
    - formatter
      - formatter.go
    - log.go
  - main.go
```

#### 3.2 formatter.go

```go
package formatter

import (
	"bytes"
	"fmt"
	"path"

	logrus "github.com/sirupsen/logrus"
)
//颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

//实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    //根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
    //自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
        //自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
        fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m  %s\n", timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}
```

#### 3.3 log.go

```go
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
	mLog := logrus.New() //新建一个实例
	mLog.SetOutput(os.Stderr) //设置输出类型
	mLog.SetReportCaller(true) //开启返回函数名和行号
	mLog.SetFormatter(&LogFormatter{}) //设置自己定义的Formatter
	mLog.SetLevel(logrus.DebugLevel) //设置最低的Level
	return &Log{
		log: mLog,
	}
}
//封装一些会用到的方法
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
```

#### 3.4 main.go

```go
package main

import (
	. "./log"
)

func Demo() {
	Logger.Info("i'm demo")

}

func main() {
	Demo()
}

//输出，其中[info]为蓝色
[2022-01-21 10:10:47] [info] entry.go:359 github.com/sirupsen/logrus.(*Entry).Logln i'm demo
```
