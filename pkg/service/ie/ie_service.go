package ieservice

import (
	"context"
	"time"

	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/database/postgresdb"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	ierepo "github.com/nhuongmh/cfvs.jpx/pkg/service/ie/repo"
	"github.com/pkg/errors"
)

type IEservice struct {
	contextTimeout time.Duration
	repo           *ierepo.IErepo
	env            *bootstrap.Env
}

func NewIEservice(timeout time.Duration, env *bootstrap.Env, db *postgresdb.DB) *IEservice {
	ies := &IEservice{
		contextTimeout: timeout,
		repo:           ierepo.NewIeRepo(db),
		env:            env,
	}
	return ies
}

func (ies *IEservice) SaveArticle(ctx context.Context, article *ie.Article) (*ie.Article, error) {
	article, err := ies.repo.Save(ctx, article)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save article")
	}

	return article, nil
}

func (ies *IEservice) GetArticle(ctx context.Context, id uint64) (*ie.Article, error) {
	article, err := ies.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article")
	}

	return article, nil
}

func (ies *IEservice) GetAllArticles(ctx context.Context, pageSize, page uint64) ([]*ie.Article, int, error) {
	articles, total, err := ies.repo.FindAll(ctx, pageSize, page)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get all articles")
	}

	return articles, total, nil
}

// find article by title
func (ies *IEservice) FindArticleByTitle(ctx context.Context, title string) ([]*ie.Article, error) {
	article, err := ies.repo.FindByTitle(ctx, title)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article by title")
	}

	return article, nil
}

// delete article by id
func (ies *IEservice) DeleteArticle(ctx context.Context, id uint64) error {
	err := ies.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete article")
	}

	return nil
}
