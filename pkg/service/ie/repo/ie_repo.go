package ierepo

import "github.com/nhuongmh/cfvs.jpx/pkg/database/postgresdb"

type IErepo struct {
	db *postgresdb.DB
}

func NewIeRepo(db *postgresdb.DB) *IErepo {
	return &IErepo{
		db: db,
	}
}
