package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmailLogsModel = (*customEmailLogsModel)(nil)

type (
	// EmailLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmailLogsModel.
	EmailLogsModel interface {
		emailLogsModel
		withSession(session sqlx.Session) EmailLogsModel
	}

	customEmailLogsModel struct {
		*defaultEmailLogsModel
	}
)

// NewEmailLogsModel returns a model for the database table.
func NewEmailLogsModel(conn sqlx.SqlConn) EmailLogsModel {
	return &customEmailLogsModel{
		defaultEmailLogsModel: newEmailLogsModel(conn),
	}
}

func (m *customEmailLogsModel) withSession(session sqlx.Session) EmailLogsModel {
	return NewEmailLogsModel(sqlx.NewSqlConnFromSession(session))
}
