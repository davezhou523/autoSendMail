package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GoogleGlModel = (*customGoogleGlModel)(nil)

type (
	// GoogleGlModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGoogleGlModel.
	GoogleGlModel interface {
		googleGlModel
		withSession(session sqlx.Session) GoogleGlModel
	}

	customGoogleGlModel struct {
		*defaultGoogleGlModel
	}
)

// NewGoogleGlModel returns a model for the database table.
func NewGoogleGlModel(conn sqlx.SqlConn) GoogleGlModel {
	return &customGoogleGlModel{
		defaultGoogleGlModel: newGoogleGlModel(conn),
	}
}

func (m *customGoogleGlModel) withSession(session sqlx.Session) GoogleGlModel {
	return NewGoogleGlModel(sqlx.NewSqlConnFromSession(session))
}
