package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	// log file name
	Filename string
	// Maximum size of a single log file
	MaxSize int
	// How many log files are saved
	MaxBackups int
	// How many days of logs are kept
	MaxDays int
	// Do you want to output the code location
	Caller bool
	// debug info warn error dpanic panic fatal
	FileLevel string
	// debug info warn error dpanic panic fatal
	ConsoleLevel string
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Helper struct {
	noCopy noCopy
	*zap.Logger

	fileLevel    zap.AtomicLevel
	consoleLevel zap.AtomicLevel
}

var emptyAtomicLevel = zap.NewAtomicLevel()

func (l *Helper) Attach(src *Helper) {
	l.Logger = src.Logger
	l.fileLevel = src.fileLevel
	l.consoleLevel = src.consoleLevel
}

func (l *Helper) Detach() {
	l.Logger = nil
	l.fileLevel = emptyAtomicLevel
	l.consoleLevel = emptyAtomicLevel
}

func (l *Helper) FileLevel() zap.AtomicLevel {
	return l.fileLevel
}

func (l *Helper) ConsoleLevel() zap.AtomicLevel {
	return l.consoleLevel
}

func NewHelper(options *Options, zapOptions ...zap.Option) *Helper {
	var cores []zapcore.Core
	fileLevel := zap.NewAtomicLevel()
	consoleLevel := zap.NewAtomicLevel()

	fileLevel = zap.NewAtomicLevel()
	if options.FileLevel == "" {
		fileLevel.SetLevel(zap.FatalLevel)
	} else if e := fileLevel.UnmarshalText([]byte(options.FileLevel)); e != nil {
		fileLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.Filename,
			MaxSize:    options.MaxSize, // megabytes
			MaxBackups: options.MaxBackups,
			MaxAge:     options.MaxDays, // days
		}),
		fileLevel,
	))

	if options.ConsoleLevel == "" {
		consoleLevel.SetLevel(zap.FatalLevel)
	} else if e := consoleLevel.UnmarshalText([]byte(options.ConsoleLevel)); e != nil {
		consoleLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		os.Stdout,
		consoleLevel,
	))

	if options.Caller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}

	return &Helper{
		Logger: zap.New(
			zapcore.NewTee(cores...),
			zapOptions...,
		),
		fileLevel:    fileLevel,
		consoleLevel: consoleLevel,
	}
}
