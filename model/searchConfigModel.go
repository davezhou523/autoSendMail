package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchConfigModel = (*customSearchConfigModel)(nil)

type (
	// SearchConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchConfigModel.
	SearchConfigModel interface {
		searchConfigModel
		withSession(session sqlx.Session) SearchConfigModel
	}

	customSearchConfigModel struct {
		*defaultSearchConfigModel
	}
)

// NewSearchConfigModel returns a model for the database table.
func NewSearchConfigModel(conn sqlx.SqlConn) SearchConfigModel {
	return &customSearchConfigModel{
		defaultSearchConfigModel: newSearchConfigModel(conn),
	}
}

func (m *customSearchConfigModel) withSession(session sqlx.Session) SearchConfigModel {
	return NewSearchConfigModel(sqlx.NewSqlConnFromSession(session))
}
