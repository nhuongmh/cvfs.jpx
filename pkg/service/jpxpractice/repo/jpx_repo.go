package repo

import (
	"context"
	"database/sql"

	"github.com/nhuongmh/cfvs.jpx/pkg/database/sqlite3"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/pkg/errors"
)

type practiceRepo struct {
	db *sqlite3.DB
}

func NewJpxPraticeRepo(db *sqlite3.DB) langfi.PracticeRepo {
	return &practiceRepo{
		db: db,
	}
}

func (rp *practiceRepo) AddCard(ctx context.Context, card *langfi.ReviewCard) error {
	query := rp.db.QueryBuilder.Insert("cards").
		Columns("front", "back", "properties", "status").
		Values(card.Front, card.Back, card.PropertiesToJson(), card.Status).
		Suffix("RETURNING id")

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	err = rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...).Scan(&card.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert card")
	}

	// also add fsrs data
	err = rp.AddFsrs(ctx, &card.FsrsData, card.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to insert fsrs data to database of card id = %v", card.ID)
	}

	return nil
}

func (rp *practiceRepo) GetCard(ctx context.Context, cardID uint64) (*langfi.ReviewCard, error) {
	query := rp.db.QueryBuilder.Select("id", "front", "back", "properties", "status").
		From("cards").
		Where("id = ?", cardID)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	card := langfi.ReviewCard{}
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan card")
	}
	card.SetPropertiesFromJson(properties)

	// get fsrs data
	fsrsData, err := rp.GetFsrs(ctx, card.ID)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to get fsrs data of card id = %v", card.ID)
		card.FsrsData = langfi.FSRSData{Card: fsrs.NewCard()}
	} else {
		card.FsrsData = *fsrsData
	}

	return &card, nil
}

func (rp *practiceRepo) UpdateCard(ctx context.Context, card *langfi.ReviewCard) error {
	query := rp.db.QueryBuilder.Update("cards").
		Where("id = ?", card.ID).
		Set("front", card.Front).
		Set("back", card.Back).
		Set("properties", card.PropertiesToJson()).
		Set("status", card.Status)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	_, err = rp.db.SqlDB.ExecContext(ctx, sqlCmd, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update card")
	}

	// also add fsrs data

	err = rp.UpdateFsrs(ctx, &card.FsrsData)
	if err != nil {
		return errors.Wrapf(err, "failed to update fsrs data to database of card id = %v", card.ID)
	}

	return nil
}

func (rp *practiceRepo) FetchReviewCard(ctx context.Context, groupID string) (*langfi.ReviewCard, error) {
	query := rp.db.QueryBuilder.Select("cards.id", "cards.front", "cards.back", "cards.properties", "cards.status").
		From("cards").
		Join("fsrs ON cards.id = fsrs.card_id").
		GroupBy("cards.id", "fsrs.card_id").
		Where("cards.status = ?", langfi.CARD_LEARN).
		OrderBy("fsrs.due").
		Limit(1)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	card := langfi.ReviewCard{}
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan card")
	}
	card.SetPropertiesFromJson(properties)

	// get fsrs data
	fsrsData, err := rp.GetFsrs(ctx, card.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get fsrs data of card id = %v", card.ID)
	}

	card.FsrsData = *fsrsData

	return &card, nil
}

func (rp *practiceRepo) FetchUnProcessCard(ctx context.Context, groupID string) (*langfi.ReviewCard, error) {
	query := rp.db.QueryBuilder.Select("id", "front", "back", "properties", "status").
		From("cards").
		Where("status = ?", langfi.CARD_NEW).
		OrderBy("created_at").
		Limit(1)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	var card langfi.ReviewCard
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoMoreDataAvailable
		}
		return nil, errors.Wrap(err, "failed to scan card")
	}
	card.SetPropertiesFromJson(properties)

	return &card, nil
}
