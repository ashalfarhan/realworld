package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			repopath := fmt.Sprintf("%s/src/github.com/ashalfarhan/", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
	Log = logger
	if config.Co.Env == "test" {
		logger.SetOutput(io.Discard)
	}
}

func New(key, value string) *logrus.Entry {
	return Log.WithField(key, value)
}

func GetCtx(ctx context.Context) *logrus.Entry {
	req := utils.GetReqID(ctx)
	return New("request_id", req)
}
