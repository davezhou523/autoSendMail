package model

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EmailProvidersModel = (*customEmailProvidersModel)(nil)

type (
	// EmailProvidersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailProvidersModel.
	EmailProvidersModel interface {
		emailProvidersModel
		withSession(session sqlx.Session) EmailProvidersModel
		FindAll(ctx context.Context, user_id int64, company_id int64) ([]*EmailProviders, error)
	}

	customEmailProvidersModel struct {
		*defaultEmailProvidersModel
	}
)

// NewEmailProvidersModel returns a model for the database table.
func NewEmailProvidersModel(conn sqlx.SqlConn) EmailProvidersModel {
	return &customEmailProvidersModel{
		defaultEmailProvidersModel: newEmailProvidersModel(conn),
	}
}

func (m *customEmailProvidersModel) withSession(session sqlx.Session) EmailProvidersModel {
	return NewEmailProvidersModel(sqlx.NewSqlConnFromSession(session))
}
func (m *customEmailProvidersModel) FindAll(ctx context.Context, user_id int64, company_id int64) ([]*EmailProviders, error) {
	selectBuilder := sq.Select("*").From(m.tableName())
	if user_id > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"user_id": user_id})
	}
	if company_id > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"company_id": company_id})
	}
	//SELECT * FROM `trade`.`email_providers` WHERE `password` !='' or oauth_client_secret !=''
	query, args, err := selectBuilder.
		Where(
			sq.Or{
				sq.NotEq{"password": ""},
				sq.NotEq{"oauth_client_secret": ""},
			},
		).
		OrderBy("id asc").
		ToSql()
	var resp []*EmailProviders
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
