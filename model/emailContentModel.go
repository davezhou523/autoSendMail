package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EmailContentModel = (*customEmailContentModel)(nil)

type (
	// EmailContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailContentModel.
	EmailContentModel interface {
		emailContentModel
		withSession(session sqlx.Session) EmailContentModel
		FindOneBySort(ctx context.Context, sort uint64) (*EmailContent, error)
	}

	customEmailContentModel struct {
		*defaultEmailContentModel
	}
)

// NewEmailContentModel returns a model for the database table.
func NewEmailContentModel(conn sqlx.SqlConn) EmailContentModel {
	return &customEmailContentModel{
		defaultEmailContentModel: newEmailContentModel(conn),
	}
}

func (m *customEmailContentModel) withSession(session sqlx.Session) EmailContentModel {
	return NewEmailContentModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultEmailContentModel) FindOneBySort(ctx context.Context, sort uint64) (*EmailContent, error) {
	var resp EmailContent
	query := fmt.Sprintf("select %s from %s where `sort` = ? limit 1", emailContentRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, sort)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
