package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmailContentModel = (*customEmailContentModel)(nil)

type (
	// EmailContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailContentModel.
	EmailContentModel interface {
		emailContentModel
		withSession(session sqlx.Session) EmailContentModel
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
