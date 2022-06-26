package logger

import (
	"context"
	"io"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	logger := logrus.New()
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
