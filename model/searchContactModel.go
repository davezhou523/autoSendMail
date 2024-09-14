package model

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SearchContactModel = (*customSearchContactModel)(nil)

type (
	// SearchContactModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchContactModel.
	SearchContactModel interface {
		searchContactModel
		FindAll(ctx context.Context, isSend uint64, category uint64, id uint64, email string, create_time string, page uint64, pageSize uint64) ([]*SearchContact, error)
		FindOneByEmail(ctx context.Context, email string) (*SearchContact, error)

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
func (m *defaultSearchContactModel) FindAll(ctx context.Context, isSend uint64, category uint64, id uint64, email string, create_time string, page uint64, pageSize uint64) ([]*SearchContact, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if isSend > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"is_send": isSend})
	}
	if category > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"category": category})
	}
	if id > 0 {
		selectBuilder = selectBuilder.Where(sq.Lt{"id": id})
	}
	if email == "notEmpty" {
		selectBuilder = selectBuilder.Where(sq.NotEq{"email": ""})
	} else if len(email) > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"email": email})
	}
	if len(create_time) > 0 {
		selectBuilder = selectBuilder.Where(sq.GtOrEq{"create_time": create_time})
	}

	selectBuilder = selectBuilder.Where(sq.Eq{"is_return": 0})

	offset := (page - 1) * pageSize
	query, args, err := selectBuilder.Offset(offset).Limit(pageSize).OrderBy("id asc").ToSql()
	println(query)
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
func (m *defaultSearchContactModel) FindOneByEmail(ctx context.Context, email string) (*SearchContact, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if len(email) > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"email": email})
	} else {
		return nil, errors.New("email is not empty")
	}
	selectBuilder = selectBuilder.Where(sq.Eq{"is_return": 0})

	query, args, err := selectBuilder.Limit(1).ToSql()

	var resp SearchContact
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
