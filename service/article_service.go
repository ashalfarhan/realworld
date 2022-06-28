package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ashalfarhan/realworld/cache/store"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/gosimple/slug"
	"github.com/matoous/go-nanoid/v2"
)

const (
	defaultSlugId = 8
)

type ArticleService struct {
	articleRepo   repository.ArticleRepository
	userRepo      repository.UserRepository
	tagsRepo      repository.ArticleTagsRepository
	favoritesRepo repository.ArticleFavoritesRepository
	commentRepo   repository.CommentRepository
	articleCache  store.ArticleStore
}

func NewArticleService(repo *repository.Repository, store *store.CacheStore) *ArticleService {
	return &ArticleService{
		repo.ArticleRepo,
		repo.UserRepo,
		repo.ArticleTagsRepo,
		repo.ArticleFavoritesRepo,
		repo.CommentRepo,
		store.ArticleStore,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, d *model.CreateArticleFields, username string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST CreateArticle dto:%+v, user:%q", d, username)
	d.Slug = s.CreateSlug(d.Title)
	a, err := s.articleRepo.InsertOne(ctx, d, username)
	if err != nil {
		log.Warnf("Cannot insert article args:%+v reason:%v", a, err)
		return nil, conduit.GeneralError
	}
	if tgs := len(d.TagList); tgs > 0 {
		tags := make([]repository.InsertArticleTagsArgs, tgs)
		for i, tag := range d.TagList {
			tags[i] = repository.InsertArticleTagsArgs{ArticleID: a.ID, TagName: tag}
			a.TagList = append(a.TagList, tag)
		}
		if err = s.tagsRepo.InsertBulk(ctx, tags); err != nil {
			log.Warnln("Cannot insert bulk tags:", err)
			return nil, conduit.GeneralError
		}
	}
	u, err := s.userRepo.FindOneByUsername(ctx, a.AuthorUsername)
	if err != nil {
		log.Warnf("Cannot find user:%q, reason:%v", a.AuthorUsername, err)
		return nil, conduit.GeneralError
	}
	a.Author = u.Profile(false) // TODO: Change following dynamically 
	return a, nil
}

func (s *ArticleService) GetArticleBySlug(ctx context.Context, username, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	if cached := s.articleCache.FindOneBySlug(ctx, slug, username); cached != nil {
		return cached, nil
	}

	ar, err := s.articleRepo.FindOneBySlug(ctx, username, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, conduit.BuildError(http.StatusNotFound, ErrNoArticleFound)
		}
		log.Warnln("Failed to get article by slug:", err)
		return nil, conduit.GeneralError
	}

	if err := s.PopulateArticleField(ctx, ar, username); err != nil {
		return nil, err
	}
	s.articleCache.SaveBySlug(ctx, slug, username, ar)
	return ar, nil
}

func (s *ArticleService) GetArticles(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	articles, err := s.articleRepo.Find(ctx, args)
	if err != nil {
		log.Warnln("Cannot find articles:", err)
		return nil, conduit.GeneralError
	}

	for _, a := range articles {
		if err := s.PopulateArticleField(ctx, a, args.Username); err != nil {
			return nil, err
		}
	}
	// TODO: Caching
	return articles, nil
}

func (s *ArticleService) GetArticlesFeed(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	articles, err := s.articleRepo.Find(ctx, args)
	if err != nil {
		log.Warnln("Cannot find feed articles:", err)
		return nil, conduit.GeneralError
	}

	for _, a := range articles {
		if err := s.PopulateArticleField(ctx, a, args.Username); err != nil {
			return nil, err
		}
	}
	// TODO: Caching
	return articles, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, slug, username string) *model.ConduitError {
	log := logger.GetCtx(ctx)
	a, err := s.GetArticleBySlug(ctx, username, slug)
	if err != nil {
		return err
	}
	if a.AuthorUsername != username {
		log.Warnf("Forbidden delete article author_username:%q, user:%q", a.AuthorUsername, username)
		return conduit.BuildError(http.StatusForbidden, ErrNotAllowedDeleteArticle)
	}
	if err := s.articleRepo.DeleteBySlug(ctx, slug); err != nil {
		log.Warnf("Failed to delete article by slug:%q, reason: %v", slug, err)
		return conduit.GeneralError
	}
	return nil
}

func (s *ArticleService) UpdateArticleBySlug(ctx context.Context, username, slug string, d *model.UpdateArticleFields) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("UpdateArticleBySlug user:%q, slug:%q, dto:%+v", username, slug, d)
	ar, err := s.GetArticleBySlug(ctx, username, slug)
	if err != nil {
		return nil, err
	}

	if ar.AuthorUsername != username {
		return nil, conduit.BuildError(http.StatusForbidden, ErrNotAllowedUpdateArticle)
	}

	// Updating title will update the slug
	if v := d.Title; v != nil {
		newSlug := s.CreateSlug(*v)
		d.Slug = &newSlug
	}

	if err := s.articleRepo.UpdateOneBySlug(ctx, d, ar); err != nil {
		log.Warnf("Cannot UpdateOneBySlug slug:%s, payload:%+v, reason: %v", slug, d, err)
		return nil, conduit.GeneralError
	}
	return ar, nil
}

func (s *ArticleService) PopulateArticleField(ctx context.Context, a *model.Article, username string) *model.ConduitError {
	tags, err := s.tagsRepo.FindArticleTagsByID(ctx, a.ID)
	if err != nil {
		return conduit.GeneralError
	}
	a.TagList = tags
	// a.Author.Following
	return nil
}

// TODO: Can be moved to utils
func (s *ArticleService) CreateSlug(title string) string {
	id, _ := gonanoid.New(defaultSlugId)
	slug := slug.Make(title)
	return fmt.Sprintf("%s-%s", slug, id)
}
