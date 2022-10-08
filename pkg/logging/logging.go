package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"sync"
)

var instance logger
var once sync.Once

type logger struct {
	*logrus.Entry
}

func Log() logger {
	once.Do(func() {
		l := logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s %s", filename, frame.Line), fmt.Sprintf("%s", frame.Function)
			},
			DisableColors: false,
			FullTimestamp: true,
		}

		l.SetOutput(os.Stdout)

		instance = logger{logrus.NewEntry(l)}
	})

	return instance
}
