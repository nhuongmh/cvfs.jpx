package repo

import (
	"github.com/nhuongmh/cfvs.jpx/pkg/database/sqlite3"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
)

type jpxRepo struct {
	db *sqlite3.DB
}

func NewJpxRepo(db *sqlite3.DB) jp.JpxRepository {
	return &jpxRepo{
		db: db,
	}
}
