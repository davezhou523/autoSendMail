package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchContactModel = (*customSearchContactModel)(nil)

type (
	// SearchContactModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchContactModel.
	SearchContactModel interface {
		searchContactModel
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
