package repo

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/database/sqlite3"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/pkg/errors"
)

type jpxRepo struct {
	db *sqlite3.DB
}

func NewJpxRepo(db *sqlite3.DB) jp.JpxGeneratorRepository {
	return &jpxRepo{
		db: db,
	}
}

func (rp *jpxRepo) AddCardProposals(ctx context.Context, cards *[]jp.CardProposal) (*[]jp.CardProposal, error) {
	for i := range *cards {
		card := (*cards)[i]
		err := rp.AddCardProposal(ctx, &card)
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert card proposal")
		}
	}

	return cards, nil
}

func (rp *jpxRepo) AddCardProposal(ctx context.Context, card *jp.CardProposal) error {
	query := rp.db.QueryBuilder.Insert("card_proposal").
		Columns("front", "back", "state", "formula_id").
		Values(card.Front, card.Back, card.State, card.FormulaID).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	err = rp.db.SqlDB.QueryRowContext(ctx, sql, args...).Scan(&card.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert card proposal")
	}

	return nil
}

func (rp *jpxRepo) ModifyCard(ctx context.Context, card *jp.CardProposal) error {
	query := rp.db.QueryBuilder.Update("card_proposal").
		Set("front", card.Front).
		Set("back", card.Back).
		Set("state", card.State).
		Set("formula_id", card.FormulaID).
		Where("id = ?", card.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	_, err = rp.db.SqlDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update card proposal")
	}

	return nil
}
