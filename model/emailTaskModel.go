package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmailTaskModel = (*customEmailTaskModel)(nil)

type (
	// EmailTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailTaskModel.
	EmailTaskModel interface {
		emailTaskModel
		withSession(session sqlx.Session) EmailTaskModel
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

func (m *customEmailTaskModel) withSession(session sqlx.Session) EmailTaskModel {
	return NewEmailTaskModel(sqlx.NewSqlConnFromSession(session))
}
