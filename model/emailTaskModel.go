package model

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ EmailTaskModel = (*customEmailTaskModel)(nil)

type (
	// EmailTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailTaskModel.
	EmailTaskModel interface {
		emailTaskModel
		WithSession(session sqlx.Session) EmailTaskModel
		FindOneBySort(ctx context.Context, id uint64, email string) (*EmailTask, error)
		FindAll(ctx context.Context, email string) ([]*EmailTask, error)
	}

	customEmailTaskModel struct {
		*defaultEmailTaskModel
	}
)

// NewEmailTaskModel returns a model for the database table.
func NewEmailTaskModel(conn sqlx.SqlConn) EmailTaskModel {
	return &customEmailTaskModel{
		defaultEmailTaskModel: newEmailTaskModel(conn),
	}
}

func (m *customEmailTaskModel) WithSession(session sqlx.Session) EmailTaskModel {
	return NewEmailTaskModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultEmailTaskModel) FindOneBySort(ctx context.Context, id uint64, email string) (*EmailTask, error) {
	selectBuilder := sq.Select("*").From(m.tableName())
	if len(strings.Trim(email, "")) > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"email": email})
	}
	if id > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"id": id})
	}
	query, args, err := selectBuilder.OrderByClause("id desc").Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var resp EmailTask
	err = m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultEmailTaskModel) FindAll(ctx context.Context, email string) ([]*EmailTask, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if len(strings.Trim(email, "")) > 0 {
		selectBuilder = selectBuilder.Where(sq.Like{"email": email})
	}

	query, args, err := selectBuilder.Limit(1000).ToSql()
	var resp []*EmailTask
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
