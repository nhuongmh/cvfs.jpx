package ierepo

import (
	"context"
	"encoding/json"

	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ir *IErepo) SaveArticleReading(ctx context.Context, articleReading *ie.ArticleReading) (*ie.ArticleReading, error) {
	questionsJSON, err := json.Marshal(articleReading.Questions)
	if err != nil {
		return nil, errors.Wrap(err, "marshal questions to JSON")
	}

	query := ir.db.QueryBuilder.Insert("article_reading").
		Columns("article_id", "questions", "article_status").
		Values(articleReading.ArticleID, questionsJSON, articleReading.Status).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&articleReading.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return articleReading, nil
}

func (ir *IErepo) FindArticleReadingByID(ctx context.Context, id uint64) (*ie.ArticleReading, error) {
	query := ir.db.QueryBuilder.Select("id", "article_id", "questions", "article_status", "created_at", "updated_at").
		From("article_reading").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var article ie.ArticleReading
	var questionsJSON []byte
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&article.ID,
		&article.ArticleID,
		&questionsJSON,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	err = json.Unmarshal(questionsJSON, &article.Questions)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal questions JSON")
	}

	return &article, nil
}

func (ir *IErepo) FindReadingByArticleId(ctx context.Context, articleId uint64) (*ie.ArticleReading, error) {
	query := ir.db.QueryBuilder.Select("id", "article_id", "questions", "article_status", "created_at", "updated_at").
		From("article_reading").
		Where("article_id = ?", articleId).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var article ie.ArticleReading
	var questionsJSON []byte
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&article.ID,
		&article.ArticleID,
		&questionsJSON,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	err = json.Unmarshal(questionsJSON, &article.Questions)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal questions JSON")
	}

	return &article, nil
}

func (ir *IErepo) DeleteArticleReading(ctx context.Context, id uint64) error {
	query := ir.db.QueryBuilder.Delete("article_reading").Where("id = ?", id)

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

func (ir *IErepo) UpdateArticleReading(ctx context.Context, articleReading *ie.ArticleReading) (*ie.ArticleReading, error) {
	questionsJSON, err := json.Marshal(articleReading.Questions)
	if err != nil {
		return nil, errors.Wrap(err, "marshal questions to JSON")
	}
	query := ir.db.QueryBuilder.Update("article_reading").
		Set("questions", questionsJSON).
		Set("article_status", articleReading.Status).
		Where("id = ?", articleReading.ID).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&articleReading.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return articleReading, nil
}

func (ir *IErepo) UpdateArticleReadingStatus(ctx context.Context, id uint64, status string) error {
	query := ir.db.QueryBuilder.Update("article_reading").
		Set("article_status", status).
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	_, err = ir.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "query row")
	}

	return nil
}
