package diary

import (
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const ErrLogOut = "log output: %v\n"

type Logger interface {
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type logger struct {
	log kitlog.Logger
}

func NewLogger(verbose bool) Logger {
	log := kitlog.NewLogfmtLogger(os.Stderr)
	if verbose {
		log = level.NewFilter(log, level.AllowDebug())
	} else {
		log = level.NewFilter(log, level.AllowAll())
	}
	return &logger{
		log: log,
	}
}

func (l logger) Debug(args ...interface{}) {
	err := level.Debug(l.log).Log(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrLogOut, err)
	}
}

func (l logger) Warn(args ...interface{}) {
	err := level.Warn(l.log).Log(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrLogOut, err)
	}
}

func (l logger) Error(args ...interface{}) {
	err := level.Error(l.log).Log(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrLogOut, err)
	}
}
