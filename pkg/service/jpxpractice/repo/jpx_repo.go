package repo

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
		Columns("front", "back", "properties", "status", "card_group").
		Values(card.Front, card.Back, card.PropertiesToJson(), card.Status, card.Group).
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
	query := rp.db.QueryBuilder.Select("id", "front", "back", "properties", "status", "card_group").
		From("cards").
		Where(sq.Eq{"id": cardID})

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	card := langfi.ReviewCard{}
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status, &card.Group)
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

func (rp *practiceRepo) GetCardByFront(ctx context.Context, front string) (*[]langfi.ReviewCard, error) {
	query := rp.db.QueryBuilder.Select("id", "front", "back", "properties", "status", "card_group").
		From("cards").
		Where(sq.Eq{"front": front})

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	rows, err := rp.db.SqlDB.QueryContext(ctx, sqlCmd, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query SQL")
	}
	defer rows.Close()
	cards := []langfi.ReviewCard{}
	for rows.Next() {
		var card langfi.ReviewCard
		var properties string
		if err := rows.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status, &card.Group); err != nil {
			return &cards, errors.Wrap(err, "failed to scan SQL")
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

		cards = append(cards, card)
	}
	if err = rows.Err(); err != nil {
		return &cards, errors.Wrap(err, "failed to scan SQL")
	}

	return &cards, nil
}

func (rp *practiceRepo) UpdateCard(ctx context.Context, card *langfi.ReviewCard) error {
	query := rp.db.QueryBuilder.Update("cards").
		Where("id = ?", card.ID).
		Set("front", card.Front).
		Set("back", card.Back).
		Set("properties", card.PropertiesToJson()).
		Set("status", card.Status).
		Set("card_group", card.Group)

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

func (rp *practiceRepo) FetchReviewCard(ctx context.Context, group string) (*langfi.ReviewCard, error) {
	logger.Log.Info().Msgf("FetchReviewCard group = %v", group)
	query := rp.db.QueryBuilder.Select("cards.id", "cards.front", "cards.back", "cards.properties", "cards.status", "cards.card_group").
		From("cards").
		Join("fsrs ON cards.id = fsrs.card_id").
		GroupBy("cards.id", "fsrs.card_id").
		Where(sq.And{sq.Eq{"cards.status": langfi.CARD_LEARN}, sq.Eq{"cards.card_group": group}}).
		OrderBy("fsrs.due").
		Limit(1)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	card := langfi.ReviewCard{}
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status, &card.Group)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoMoreDataAvailable
		}
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

func (rp *practiceRepo) FetchUnProcessCard(ctx context.Context, group string) (*langfi.ReviewCard, error) {

	query := rp.db.QueryBuilder.Select("id", "front", "back", "properties", "status", "card_group").
		From("cards").
		Where(sq.And{sq.Eq{"status": langfi.CARD_NEW}, sq.Eq{"card_group": group}}).
		OrderBy("created_at").
		Limit(1)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}
	logger.Log.Info().Msgf("FetchUnProcessCard group = %v, sql=`%v`, args=%v", group, sqlCmd, args)

	row := rp.db.SqlDB.QueryRowContext(ctx, sqlCmd, args...)
	var card langfi.ReviewCard
	var properties string
	err = row.Scan(&card.ID, &card.Front, &card.Back, &properties, &card.Status, &card.Group)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoMoreDataAvailable
		}
		return nil, errors.Wrap(err, "failed to scan card")
	}
	card.SetPropertiesFromJson(properties)

	return &card, nil
}

func (rp *practiceRepo) DeleteNewCard(ctx context.Context) error {
	logger.Log.Info().Msg("DeleteNewCard")
	query := rp.db.QueryBuilder.Delete("cards").Where("status = ?", langfi.CARD_NEW)

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}
	_, err = rp.db.SqlDB.ExecContext(ctx, sqlCmd, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete NEW card")
	}
	return nil
}

func (rp *practiceRepo) GetGroupStats(ctx context.Context) (*[]langfi.GroupSummaryDto, error) {
	query := rp.db.QueryBuilder.Select("card_group", "count(*) as num_cards",
		fmt.Sprintf("COUNT(CASE WHEN status = '%s' THEN 1 END) AS card_new", langfi.CARD_NEW),
		fmt.Sprintf("COUNT(CASE WHEN status = '%s' THEN 1 END) AS card_learn", langfi.CARD_LEARN),
		fmt.Sprintf("COUNT(CASE WHEN status = '%s' THEN 1 END) AS card_discard", langfi.CARD_DISCARD),
		fmt.Sprintf("COUNT(CASE WHEN status = '%s' THEN 1 END) AS card_save", langfi.CARD_SAVE)).
		From("cards").
		GroupBy("card_group")

	sqlCmd, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	rows, err := rp.db.SqlDB.QueryContext(ctx, sqlCmd, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query SQL")
	}
	defer rows.Close()
	groups := []langfi.GroupSummaryDto{}
	for rows.Next() {
		var group langfi.GroupSummaryDto
		if err := rows.Scan(&group.Group, &group.NumCards, &group.Proposal, &group.Learning, &group.Discard, &group.Save); err != nil {
			return &groups, errors.Wrap(err, "failed to scan SQL")
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return &groups, errors.Wrap(err, "failed to scan SQL")
	}

	return &groups, nil
}
