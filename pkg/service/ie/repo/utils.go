package ierepo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (ir *IErepo) getTotalRowsOfQuery(ctx context.Context, sub squirrel.SelectBuilder) (int, error) {
	sql, args, err := sub.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build query")
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM (%s)", sql)
	var totalRows uint64
	err = ir.db.QueryRow(ctx, query, args...).Scan(&totalRows)
	if err != nil {
		return 0, errors.Wrap(err, "query")
	}
	return int(totalRows), nil
}
