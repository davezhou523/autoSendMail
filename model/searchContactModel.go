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
		FindAll(ctx context.Context, user_id int64, company_id int64, category int64, id uint64, email string, create_time string, page uint64, pageSize uint64, contentId uint64) ([]*SearchContact, error)
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
func (m *defaultSearchContactModel) FindAll(ctx context.Context, user_id int64, company_id int64, category int64, id uint64, email string, create_time string, page uint64, pageSize uint64, contentId uint64) ([]*SearchContact, error) {
	// 定义子查询 SQL 语句
	subQuery := sq.Select("email").
		From("email_task").
		Where(sq.Eq{"content_id": contentId})

	// 构建主查询
	selectBuilder := sq.Select("*").From(m.tableName())
	if len(create_time) > 0 {
		selectBuilder = selectBuilder.Where(sq.GtOrEq{"create_time": create_time})
	}
	if len(email) > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"email": email})
	}
	if user_id > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"user_id": user_id})
	}
	if company_id > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"company_id": company_id})
	}

	if category > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"category": category})
	}
	if id > 0 {
		selectBuilder = selectBuilder.Where(sq.Lt{"id": id})
	}

	offset := (page - 1) * pageSize
	query, args, err := selectBuilder.
		Where(sq.Expr("email NOT IN (?)", subQuery)).
		Where(sq.NotEq{"email": ""}).
		Where(sq.Eq{"is_send": 1}).
		Offset(offset).Limit(pageSize).
		OrderBy("id asc").
		ToSql()

	//fmt.Println(query)
	//fmt.Println(args)

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
