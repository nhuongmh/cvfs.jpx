package repo

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/pkg/errors"
)

func (rp *practiceRepo) AddFsrs(ctx context.Context, fsrsd *langfi.FSRSData, cardID uint64) error {
	query := rp.db.QueryBuilder.Insert("fsrs").
		Columns("card_id", "due", "stability", "difficulty", "elapsed_days", "scheduled_days",
			"reps", "lapses", "state", "last_review").
		Values(cardID, fsrsd.Due, fsrsd.Stability, fsrsd.Difficulty, fsrsd.ElapsedDays,
			fsrsd.ScheduledDays, fsrsd.Reps, fsrsd.Lapses, fsrsd.State, fsrsd.LastReview).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	err = rp.db.SqlDB.QueryRowContext(ctx, sql, args...).Scan(&fsrsd.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert fsrs")
	}

	return nil
}

func (rp *practiceRepo) GetFsrs(ctx context.Context, cardID uint64) (*langfi.FSRSData, error) {
	query := rp.db.QueryBuilder.Select("id", "due", "stability", "difficulty", "elapsed_days", "scheduled_days",
		"reps", "lapses", "state", "last_review").
		From("fsrs").
		Where("card_id = ?", cardID)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	row := rp.db.SqlDB.QueryRowContext(ctx, sql, args...)
	fsrsData := langfi.FSRSData{}

	err = row.Scan(&fsrsData.ID, &fsrsData.Due, &fsrsData.Stability, &fsrsData.Difficulty, &fsrsData.ElapsedDays,
		&fsrsData.ScheduledDays, &fsrsData.Reps, &fsrsData.Lapses, &fsrsData.State, &fsrsData.LastReview)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan fsrs")
	}

	return &fsrsData, nil
}

func (rp *practiceRepo) UpdateFsrs(ctx context.Context, fsrsd *langfi.FSRSData) error {
	query := rp.db.QueryBuilder.Update("fsrs").
		Where("id = ?", fsrsd.ID).
		Set("due", fsrsd.Due).
		Set("stability", fsrsd.Stability).
		Set("difficulty", fsrsd.Difficulty).
		Set("elapsed_days", fsrsd.ElapsedDays).
		Set("scheduled_days", fsrsd.ScheduledDays).
		Set("reps", fsrsd.Reps).
		Set("lapses", fsrsd.Lapses).
		Set("state", fsrsd.State).
		Set("last_review", fsrsd.LastReview)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	_, err = rp.db.SqlDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update fsrs")
	}

	return nil
}
