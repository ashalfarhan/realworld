package conduit

import (
	"io"
	"os"

	"github.com/ashalfarhan/realworld/config"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func NewLogger(key, value string) *logrus.Entry {
	var out io.Writer = os.Stdout
	if config.Co.Env == "test" {
		out = io.Discard
	}

	l := logrus.WithField(key, value)
	l.Logger.SetOutput(out)
	return l
}
