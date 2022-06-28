package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *ArticleService) GetAllTags(ctx context.Context) ([]string, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	tags, err := s.tagsRepo.FindAllTags(ctx)
	if err != nil {
		log.Warnln("Cannot find all tags reason:", err)
		return nil, conduit.BuildError(http.StatusInternalServerError, nil)
	}
	return tags, nil
}
