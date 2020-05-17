package diary

import (
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

const ErrLogOut = "log output: %v\n"

type Logger interface {
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type logger struct {
	log     kitlog.Logger
	logErr  kitlog.Logger
	verbose bool
}

func NewLogger(verbose bool) Logger {
	return &logger{
		log:     kitlog.NewLogfmtLogger(os.Stdout),
		logErr:  kitlog.NewLogfmtLogger(os.Stderr),
		verbose: verbose,
	}
}

func (l logger) Debug(args ...interface{}) {
	if l.verbose {
		err := l.log.Log(
			interfaceSlice(
				"level", "debug",
				args...,
			)...,
		)
		if err != nil {
			fmt.Printf(ErrLogOut, err)
		}
	}
}

func (l logger) Warn(args ...interface{}) {
	if l.verbose {
		err := l.log.Log(
			interfaceSlice(
				"level", "warn",
				args...,
			)...,
		)
		if err != nil {
			fmt.Printf(ErrLogOut, err)
		}
	}
}

func (l logger) Error(args ...interface{}) {
	err := l.logErr.Log(
		interfaceSlice(
			"level", "error",
			args...,
		)...,
	)
	if err != nil {
		fmt.Printf(ErrLogOut, err)
	}
}

func interfaceSlice(level, name string, args ...interface{}) []interface{} {
	interfaceSlice := make([]interface{}, len(args)+2)
	interfaceSlice[0] = level
	interfaceSlice[1] = name
	_ = copy(interfaceSlice[2:], args)
	return interfaceSlice
}
