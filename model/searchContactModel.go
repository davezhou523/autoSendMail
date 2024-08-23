package model

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SearchContactModel = (*customSearchContactModel)(nil)

type (
	// SearchContactModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchContactModel.
	SearchContactModel interface {
		searchContactModel
		FindAll(ctx context.Context, isSend uint64, category uint64) ([]*SearchContact, error)

		withSession(session sqlx.Session) SearchContactModel
	}

	customSearchContactModel struct {
		*defaultSearchContactModel
	}
)

// NewSearchContactModel returns a model for the database table.
func NewSearchContactModel(conn sqlx.SqlConn) SearchContactModel {
	return &customSearchContactModel{
		defaultSearchContactModel: newSearchContactModel(conn),
	}
}

func (m *customSearchContactModel) withSession(session sqlx.Session) SearchContactModel {
	return NewSearchContactModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultSearchContactModel) FindAll(ctx context.Context, isSend uint64, category uint64) ([]*SearchContact, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if isSend > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"is_send": isSend})
	}
	if category > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"category": category})
	}
	query, args, err := selectBuilder.Limit(1000).ToSql()
	var resp []*SearchContact
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
