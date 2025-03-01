package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchConfigRunModel = (*customSearchConfigRunModel)(nil)

type (
	// SearchConfigRunModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchConfigRunModel.
	SearchConfigRunModel interface {
		searchConfigRunModel
		withSession(session sqlx.Session) SearchConfigRunModel
	}

	customSearchConfigRunModel struct {
		*defaultSearchConfigRunModel
	}
)

// NewSearchConfigRunModel returns a model for the database table.
func NewSearchConfigRunModel(conn sqlx.SqlConn) SearchConfigRunModel {
	return &customSearchConfigRunModel{
		defaultSearchConfigRunModel: newSearchConfigRunModel(conn),
	}
}

func (m *customSearchConfigRunModel) withSession(session sqlx.Session) SearchConfigRunModel {
	return NewSearchConfigRunModel(sqlx.NewSqlConnFromSession(session))
}
