package repo

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/database/sqlite3"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
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
	return nil, model.ErrNotImplemented
}
func (rp *jpxRepo) ModifyCard(ctx context.Context, card *jp.CardProposal) error {
	return model.ErrNotImplemented
}
