package ierepo

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ir *IErepo) SaveVocabList(ctx context.Context, list *ie.IeVocabList) (*ie.IeVocabList, error) {
	query := ir.db.QueryBuilder.Insert("ie_vocab_list").
		Columns("name", "article_id").
		Values(list.Name, list.RefArticleID).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&list.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	//save all vocabs of this list
	for _, word := range list.Vocabs {
		word.VocabListId = list.ID
		_, err = ir.SaveVocab(ctx, &word)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("SaveVocabList: save vocab failed")
		}
	}

	return list, nil
}

func (ir *IErepo) GetAllVocabList(ctx context.Context, fetchVocabs bool, limit, skip uint64) (*[]ie.IeVocabList, int, error) {
	query := ir.db.QueryBuilder.Select("id", "name", "article_id", "created_at", "updated_at").
		From("ie_vocab_list")

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
		return nil, 0, errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	var lists []ie.IeVocabList
	for rows.Next() {
		var list ie.IeVocabList
		err = rows.Scan(
			&list.ID,
			&list.Name,
			&list.RefArticleID,
			&list.CreatedAt,
			&list.UpdatedAt)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scan row")
		}
		if fetchVocabs {
			words, err := ir.FindVocabByListId(ctx, list.ID)
			if err != nil {
				return nil, 0, errors.Wrap(err, "find vocab by list id")
			}
			list.Vocabs = *words
		}
		lists = append(lists, list)
	}

	total, err := ir.getTotalRowsOfQuery(ctx, query)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get total rows")
		total = len(lists)
	}

	return &lists, total, nil
}

func (ir *IErepo) GetAllVocabListByArticleId(ctx context.Context, articleId uint64) (*ie.IeVocabList, error) {
	query := ir.db.QueryBuilder.Select("id", "name", "article_id", "created_at", "updated_at").
		From("ie_vocab_list").
		Where("article_id = ?", articleId)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var list ie.IeVocabList
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&list.ID,
		&list.Name,
		&list.RefArticleID,
		&list.CreatedAt,
		&list.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	vocabs, err := ir.FindVocabByListId(ctx, list.ID)
	if err != nil {
		return nil, errors.Wrap(err, "find vocab by list id")
	}
	list.Vocabs = *vocabs
	return &list, nil
}

func (ir *IErepo) FindVocabListByID(ctx context.Context, id uint64, fetchVocab bool) (*ie.IeVocabList, error) {
	query := ir.db.QueryBuilder.Select("id", "name", "article_id", "created_at", "updated_at").
		From("ie_vocab_list").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var list ie.IeVocabList
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&list.ID,
		&list.Name,
		&list.RefArticleID,
		&list.CreatedAt,
		&list.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	if fetchVocab {
		vocabs, err := ir.FindVocabByListId(ctx, list.ID)
		if err != nil {
			return nil, errors.Wrap(err, "find vocab by list id")
		}
		list.Vocabs = *vocabs
	}
	return &list, nil
}

func (ir *IErepo) DeleteVocabList(ctx context.Context, id uint64) error {
	query := ir.db.QueryBuilder.Delete("ie_vocab_list").Where("id = ?", id)

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

func (ir *IErepo) UpdateVocabList(ctx context.Context, list *ie.IeVocabList) (*ie.IeVocabList, error) {
	query := ir.db.QueryBuilder.Update("ie_vocab_list").
		Set("name", list.Name).
		Set("article_id", list.RefArticleID).
		Where("id = ?", list.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	_, err = ir.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec")
	}

	// update all vocabs of this list
	for _, word := range list.Vocabs {
		savedVocab, err := ir.FindOneVocabWordInList(ctx, word.Word, list.ID)
		word.VocabListId = list.ID
		if err != nil {
			_, err = ir.SaveVocab(ctx, &word)
			if err != nil {
				logger.Log.Warn().Err(err).Msg("UpdateVocabList: save vocab failed")
			}
		} else {
			word.ID = savedVocab.ID
			_, err = ir.UpdateVocab(ctx, &word)
			if err != nil {
				logger.Log.Warn().Err(err).Msg("UpdateVocabList: update vocab failed")
			}
		}
	}

	return list, nil
}
