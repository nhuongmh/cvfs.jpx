package ierepo

import (
	"context"
	"encoding/json"

	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ir *IErepo) SaveTestSubmission(ctx context.Context, testResult *ie.TestResult) (*ie.TestResult, error) {
	questionsJSON, err := json.Marshal(testResult.QuestionResults)
	if err != nil {
		return nil, errors.Wrap(err, "marshal questions result to JSON")
	}

	query := ir.db.QueryBuilder.Insert("article_test_result").
		Columns("article_reading_id", "questions_result", "score").
		Values(testResult.ArticleReadingId, questionsJSON, testResult.Score).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&testResult.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return testResult, nil
}

func (ir *IErepo) GetTestSubmissionById(ctx context.Context, id uint64) (*ie.TestResult, error) {
	query := ir.db.QueryBuilder.Select("id", "article_reading_id", "questions_result", "score", "created_at", "updated_at").
		From("article_test_result").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var testResult ie.TestResult
	var questionsJSON []byte
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&testResult.ID,
		&testResult.ArticleReadingId,
		&questionsJSON,
		&testResult.Score,
		&testResult.CreatedAt,
		&testResult.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	err = json.Unmarshal(questionsJSON, &testResult.QuestionResults)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal questions JSON")
	}
	return &testResult, nil
}

func (ir *IErepo) FindSubmissionByReadingId(ctx context.Context, readingId uint64) (*[]ie.TestResult, error) {
	query := ir.db.QueryBuilder.Select("id", "article_reading_id", "questions_result", "score", "created_at", "updated_at").
		From("article_test_result").
		Where("article_reading_id = ?", readingId)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	rows, err := ir.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var testResults []ie.TestResult
	for rows.Next() {
		var testResult ie.TestResult
		var questionsJSON []byte
		err = rows.Scan(
			&testResult.ID,
			&testResult.ArticleReadingId,
			&questionsJSON,
			&testResult.Score,
			&testResult.CreatedAt,
			&testResult.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan row")
		}

		err = json.Unmarshal(questionsJSON, &testResult.QuestionResults)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal questions JSON")
		}
	}

	return &testResults, nil
}

func (ir *IErepo) DeleteSubmission(ctx context.Context, id uint64) error {
	query := ir.db.QueryBuilder.Delete("article_test_result").Where("id = ?", id)

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
