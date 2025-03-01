package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GoogleLrModel = (*customGoogleLrModel)(nil)

type (
	// GoogleLrModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGoogleLrModel.
	GoogleLrModel interface {
		googleLrModel
		withSession(session sqlx.Session) GoogleLrModel
	}

	customGoogleLrModel struct {
		*defaultGoogleLrModel
	}
)

// NewGoogleLrModel returns a model for the database table.
func NewGoogleLrModel(conn sqlx.SqlConn) GoogleLrModel {
	return &customGoogleLrModel{
		defaultGoogleLrModel: newGoogleLrModel(conn),
	}
}

func (m *customGoogleLrModel) withSession(session sqlx.Session) GoogleLrModel {
	return NewGoogleLrModel(sqlx.NewSqlConnFromSession(session))
}
