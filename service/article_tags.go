package service

import (
	"context"
	"net/http"
)

func (s *ArticleService) GetAllTags(ctx context.Context) ([]string, *ServiceError) {
	tags, err := s.tagsRepo.FindAllTags(ctx)
	if err != nil {
		s.logger.Printf("Cannot FindAllTags, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return tags, nil
}
