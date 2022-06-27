package logger

import (
	"context"
	"io"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/sirupsen/logrus"
)

func Configure() {
	if config.Env == "test" {
		logrus.SetOutput(io.Discard)
	}
}

func GetCtx(ctx context.Context) *logrus.Entry {
	return logrus.WithField("request_id", utils.GetReqID(ctx))
}
