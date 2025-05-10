package ierepo

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ir *IErepo) Save(ctx context.Context, article *ie.Article) (*ie.Article, error) {
	query := ir.db.QueryBuilder.Insert("ie_articles").
		Columns("title", "content", "origin", "author", "cover_image", "publish_date").
		Values(article.Title, article.Content, article.Origin, article.Author, article.Image, article.PublishDate).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&article.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return article, nil
}

func (ir *IErepo) FindByID(ctx context.Context, id uint64) (*ie.Article, error) {
	query := ir.db.QueryBuilder.Select("id", "title", "content", "origin", "author", "cover_image", "publish_date", "created_at", "updated_at").
		From("ie_articles").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var article ie.Article
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Origin,
		&article.Author,
		&article.Image,
		&article.PublishDate,
		&article.CreatedAt,
		&article.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return &article, nil
}

// omit content
func (ir *IErepo) FindAll(ctx context.Context, limit, skip uint64) ([]*ie.Article, int, error) {
	query := ir.db.QueryBuilder.Select("id", "title", "origin", "author", "cover_image", "publish_date", "created_at", "updated_at").
		From("ie_articles")

	pageQuery := query.
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := pageQuery.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "build query")
	}

	logger.Log.Debug().Str("sql find all", sql).Msg("query")

	rows, err := ir.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "query")
	}
	defer rows.Close()

	articles := []*ie.Article{}
	for rows.Next() {
		var article ie.Article
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Origin,
			&article.Author,
			&article.Image,
			&article.PublishDate,
			&article.CreatedAt,
			&article.UpdatedAt)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scan row")
		}
		articles = append(articles, &article)
	}

	total, err := ir.getTotalRowsOfQuery(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get total rows")
		total = len(articles)
	}

	return articles, total, nil
}

func (ir *IErepo) Delete(ctx context.Context, id uint64) error {
	query := ir.db.QueryBuilder.Delete("ie_articles").Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	_, err = ir.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "exec")
	}

	return nil
}

func (ir *IErepo) Update(ctx context.Context, article *ie.Article) (*ie.Article, error) {
	query := ir.db.QueryBuilder.Update("ie_articles").
		Set("title", article.Title).
		Set("content", article.Content).
		Set("origin", article.Origin).
		Set("author", article.Author).
		Set("cover_image", article.Image).
		Set("publish_date", article.PublishDate).
		Where("id = ?", article.ID).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&article.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return article, nil
}

// omit content
func (ir *IErepo) FindByTitle(ctx context.Context, title string) ([]*ie.Article, error) {
	query := ir.db.QueryBuilder.Select("id", "title", "origin", "author", "cover_image", "publish_date", "created_at", "updated_at").
		From("ie_articles").
		Where("title = ?", title)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	rows, err := ir.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var articles []*ie.Article
	for rows.Next() {
		var article ie.Article
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Origin,
			&article.Author,
			&article.Image,
			&article.PublishDate,
			&article.CreatedAt,
			&article.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan row")
		}
		articles = append(articles, &article)
	}

	return articles, nil
}
