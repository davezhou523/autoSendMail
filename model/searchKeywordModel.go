package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchKeywordModel = (*customSearchKeywordModel)(nil)

type (
	// SearchKeywordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchKeywordModel.
	SearchKeywordModel interface {
		searchKeywordModel
		withSession(session sqlx.Session) SearchKeywordModel
	}

	customSearchKeywordModel struct {
		*defaultSearchKeywordModel
	}
)

// NewSearchKeywordModel returns a model for the database table.
func NewSearchKeywordModel(conn sqlx.SqlConn) SearchKeywordModel {
	return &customSearchKeywordModel{
		defaultSearchKeywordModel: newSearchKeywordModel(conn),
	}
}

func (m *customSearchKeywordModel) withSession(session sqlx.Session) SearchKeywordModel {
	return NewSearchKeywordModel(sqlx.NewSqlConnFromSession(session))
}
