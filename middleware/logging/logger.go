package logging

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

type OupputType string

var (
	logPath    = "log/backend.log"
	project    = "TwoPartyEddsa"
	size       = 20
	age        = 7
	backups    = 10
	JsonFormat = "json"
	TextFormat = "text"

	ONLY_TERMINAL OupputType = "only terminal"
	ONLY_FILE     OupputType = "only file"
	MULT_OUTPUT   OupputType = "output to file and terminal"
)

type logOptions struct {
	logPath string     // 日志路径
	project string     // 项目名称
	size    int        // 单个日志文件大小
	age     int        // 日志文件保留时间
	backups int        // 备份
	output  OupputType // 输出类型
	format  string     // log format
}

func newLogOps() *logOptions {
	return &logOptions{
		logPath: logPath,
		project: project,
		size:    size,
		age:     age,
		backups: backups,
		output:  MULT_OUTPUT,
		format:  JsonFormat,
	}
}

type LogOps func(this *logOptions)

func WithLogPath(path string) LogOps {
	return func(this *logOptions) {
		this.logPath = path
	}
}

func WithLogProject(projectName string) LogOps {
	return func(this *logOptions) {
		this.project = project
	}
}

func WithLogSize(size int) LogOps {
	return func(this *logOptions) {
		this.size = size
	}
}

func WithLogAge(age int) LogOps {
	return func(this *logOptions) {
		this.age = age
	}
}

func WithLogBackups(backups int) LogOps {
	return func(this *logOptions) {
		this.backups = backups
	}
}

func WithOutput(outputType OupputType) LogOps {
	return func(this *logOptions) {
		this.output = outputType
	}
}

func WithFormat(format string) LogOps {
	if format != JsonFormat && format != TextFormat {
		panic("invalid format")
	}
	return func(this *logOptions) {
		this.format = format
	}
}

// 初始化日志
func InitLogger(ops ...LogOps) *logrus.Entry {
	logops := newLogOps()

	for _, f := range ops {
		f(logops)
	}

	log := logrus.New()

	logger := &lumberjack.Logger{
		Filename:   logops.logPath,
		MaxSize:    logops.size,
		MaxAge:     logops.age,
		MaxBackups: logops.backups,
		LocalTime:  true,
		Compress:   true,
	}

	// set output
	switch logops.output {
	case ONLY_TERMINAL:
		output := io.MultiWriter(os.Stdout)
		log.SetOutput(output)
	case ONLY_FILE:
		log.SetOutput(logger)
	case MULT_OUTPUT:
		output := io.MultiWriter(logger, os.Stdout)
		log.SetOutput(output)
	default:
		output := io.MultiWriter(logger, os.Stdout)
		log.SetOutput(output)
	}

	// set output format
	if logops.format == JsonFormat {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	newLog := log.WithField("psm", logops.project)
	return newLog
}
