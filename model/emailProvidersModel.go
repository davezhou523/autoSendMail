package model

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ EmailProvidersModel = (*customEmailProvidersModel)(nil)

type (
	// EmailProvidersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailProvidersModel.
	EmailProvidersModel interface {
		emailProvidersModel
		WithSession(session sqlx.Session) EmailProvidersModel
		FindAll(ctx context.Context, user_id int64, company_id int64) ([]*EmailProviders, error)
		ResetDailyCount()
		IncrementSent(ctx context.Context, id int64) (int64, error)
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

func (m *customEmailProvidersModel) WithSession(session sqlx.Session) EmailProvidersModel {
	return NewEmailProvidersModel(sqlx.NewSqlConnFromSession(session))
}

// 重置每日限额
func (m *customEmailProvidersModel) ResetDailyCount() {

	query := `UPDATE email_providers SET sent_count = 0,reset_time=?`
	now := time.Now()
	// 计算明天日期（+1天）
	tomorrow := now.AddDate(0, 0, 1)
	// 格式化为 YYYY-MM-DD
	dateStr := tomorrow.Format("2006-01-02")
	_, err := m.conn.Exec(query, dateStr)
	fmt.Println(err)

}

// 原子增加发送计数
func (m *customEmailProvidersModel) IncrementSent(ctx context.Context, id int64) (int64, error) {
	query := `UPDATE email_providers SET sent_count = sent_count + 1 WHERE id = ? AND sent_count < daily_limit`
	result, err := m.conn.ExecCtx(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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
