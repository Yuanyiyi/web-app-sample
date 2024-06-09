package runtime

import (
	"fmt"
	joonix "github.com/joonix/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const sourceKey = "source"

// stackTracer is the pkg/errors stacktrace interface
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// replace the standard glog error logger, with a logrus one
func init() {
	logrus.SetFormatter(&joonix.FluentdFormatter{})
}

// SetLevel select level to filter logger output
func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

// HandleError wraps runtime.HandleError so that it is possible to
// use WithField with logrus.
func HandleError(logger *logrus.Entry, err error) {
	if logger != nil {
		// it's a bit of a double handle, but I can't see a better way to do it
		logger.WithError(err).Error()
	}
}

// Must panics if there is an error
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// NewLoggerWithSource returns a logrus.Entry to use when you want to specify an source
func NewLoggerWithSource(source string) (log *logrus.Entry) {
	return logrus.WithField(sourceKey, source)
}

// NewLoggerWithType returns a logrus.Entry to use when you want to use a data type as the source
// such as when you have a struct with methods
func NewLoggerWithType(obj interface{}) *logrus.Entry {
	return NewLoggerWithSource(fmt.Sprintf("%T", obj))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
