package ierepo

import (
	"context"
	"encoding/json"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ir *IErepo) SaveVocab(ctx context.Context, word *ie.IeVocab) (*ie.IeVocab, error) {
	pronJSON, err := json.Marshal(word.Pronunciation)
	if err != nil {
		return nil, errors.Wrap(err, "marshal pronunciation to JSON")
	}
	defJSON, err := json.Marshal(word.Definitions)
	if err != nil {
		return nil, errors.Wrap(err, "marshal definitions to JSON")
	}
	propertiesJSON, err := json.Marshal(word.Properties)
	if err != nil {
		return nil, errors.Wrap(err, "marshal properties to JSON")
	}

	query := ir.db.QueryBuilder.Insert("ie_vocab").
		Columns("vocab", "vocab_list_id", "prons", "defs", "props", "context", "freq").
		Values(word.Word, word.VocabListId, pronJSON, defJSON, propertiesJSON, word.Context, word.WordFreq).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&word.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	return word, nil
}

func (ir *IErepo) FindVocabByID(ctx context.Context, id uint64) (*ie.IeVocab, error) {
	query := ir.db.QueryBuilder.Select("id", "vocab", "vocab_list_id", "prons", "defs", "props", "context", "freq", "created_at", "updated_at").
		From("ie_vocab").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var word ie.IeVocab
	var pronJSON []byte
	var defJSON []byte
	var propsJSON []byte
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&word.ID,
		&word.Word,
		&pronJSON,
		&defJSON,
		&propsJSON,
		&word.Context,
		&word.WordFreq,
		&word.CreatedAt,
		&word.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	err = json.Unmarshal(pronJSON, &word.Pronunciation)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal pronunciation JSON")
	}
	err = json.Unmarshal(defJSON, &word.Definitions)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal definitions JSON")
	}
	err = json.Unmarshal(propsJSON, &word.Properties)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal properties JSON")
	}

	return &word, nil
}

func (ir *IErepo) FindVocabWord(ctx context.Context, searchKey string) (*[]ie.IeVocab, error) {
	query := ir.db.QueryBuilder.Select("id", "vocab", "vocab_list_id", "prons", "defs", "props", "context", "freq", "created_at", "updated_at").
		From("ie_vocab").
		Where(sq.Like{"LOWER(vocab)": strings.ToLower(searchKey)})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	rows, err := ir.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	var words []ie.IeVocab
	for rows.Next() {
		var word ie.IeVocab
		var pronJSON []byte
		var defJSON []byte
		var propsJSON []byte
		err = rows.Scan(
			&word.ID,
			&word.Word,
			&word.VocabListId,
			&pronJSON,
			&defJSON,
			&propsJSON,
			&word.Context,
			&word.WordFreq,
			&word.CreatedAt,
			&word.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan row")
		}
		err = json.Unmarshal(pronJSON, &word.Pronunciation)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal pronunciation JSON")
		}
		err = json.Unmarshal(defJSON, &word.Definitions)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal definitions JSON")
		}
		err = json.Unmarshal(propsJSON, &word.Properties)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal properties JSON")
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows iteration")
	}

	return &words, nil
}

func (ir *IErepo) FindOneVocabWordInList(ctx context.Context, searchKey string, vocabListId uint64) (*ie.IeVocab, error) {
	query := ir.db.QueryBuilder.Select("id", "vocab", "vocab_list_id", "prons", "defs", "props", "context", "freq", "created_at", "updated_at").
		From("ie_vocab").
		Where(sq.And{sq.Like{"LOWER(vocab)": strings.ToLower(searchKey)}, sq.Eq{"vocab_list_id": vocabListId}}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var word ie.IeVocab
	var pronJSON []byte
	var defJSON []byte
	var propsJSON []byte
	err = ir.db.QueryRow(ctx, sql, args...).Scan(
		&word.ID,
		&word.Word,
		&pronJSON,
		&defJSON,
		&propsJSON,
		&word.Context,
		&word.WordFreq,
		&word.CreatedAt,
		&word.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	err = json.Unmarshal(pronJSON, &word.Pronunciation)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal pronunciation JSON")
	}
	err = json.Unmarshal(defJSON, &word.Definitions)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal definitions JSON")
	}
	err = json.Unmarshal(propsJSON, &word.Properties)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal properties JSON")
	}

	return &word, nil
}

func (ir *IErepo) DeleteVocab(ctx context.Context, id uint64) error {
	query := ir.db.QueryBuilder.Delete("ie_vocab").
		Where("id = ?", id)

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

func (ir *IErepo) UpdateVocab(ctx context.Context, word *ie.IeVocab) (*ie.IeVocab, error) {
	pronJSON, err := json.Marshal(word.Pronunciation)
	if err != nil {
		return nil, errors.Wrap(err, "marshal pronunciation to JSON")
	}
	defJSON, err := json.Marshal(word.Definitions)
	if err != nil {
		return nil, errors.Wrap(err, "marshal definitions to JSON")
	}
	propertiesJSON, err := json.Marshal(word.Properties)
	if err != nil {
		return nil, errors.Wrap(err, "marshal properties to JSON")
	}

	query := ir.db.QueryBuilder.Update("ie_vocab").
		Set("vocab", word.Word).
		Set("vocab_list_id", word.VocabListId).
		Set("prons", pronJSON).
		Set("defs", defJSON).
		Set("props", propertiesJSON).
		Set("context", word.Context).
		Set("freq", word.WordFreq).
		Set("updated_at", word.UpdatedAt).
		Where("id = ?", word.ID).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = ir.db.QueryRow(ctx, sql, args...).Scan(&word.ID)
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	return word, nil
}

func (ir *IErepo) FindVocabByListId(ctx context.Context, listId uint64) (*[]ie.IeVocab, error) {
	query := ir.db.QueryBuilder.Select("id", "vocab", "vocab_list_id", "prons", "defs", "props", "context", "freq", "created_at", "updated_at").
		From("ie_vocab").
		Where("vocab_list_id = ?", listId)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	rows, err := ir.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	var words []ie.IeVocab
	for rows.Next() {
		var word ie.IeVocab
		var pronJSON []byte
		var defJSON []byte
		var propsJSON []byte
		err = rows.Scan(
			&word.ID,
			&word.Word,
			&word.VocabListId,
			&pronJSON,
			&defJSON,
			&propsJSON,
			&word.Context,
			&word.WordFreq,
			&word.CreatedAt,
			&word.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan row")
		}
		if word.Word == "" {
			continue // Skip empty words.
		}

		err = json.Unmarshal(pronJSON, &word.Pronunciation)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal pronunciation JSON")
		}
		err = json.Unmarshal(defJSON, &word.Definitions)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal definitions JSON")
		}
		err = json.Unmarshal(propsJSON, &word.Properties)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal properties JSON")
		}
		word.VocabListId = listId

		words = append(words, word)
	}
	return &words, nil
}
