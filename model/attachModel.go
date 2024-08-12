package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AttachModel = (*customAttachModel)(nil)

type (
	// AttachModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAttachModel.
	AttachModel interface {
		attachModel
		withSession(session sqlx.Session) AttachModel
	}

	customAttachModel struct {
		*defaultAttachModel
	}
)

// NewAttachModel returns a model for the database table.
func NewAttachModel(conn sqlx.SqlConn) AttachModel {
	return &customAttachModel{
		defaultAttachModel: newAttachModel(conn),
	}
}

func (m *customAttachModel) withSession(session sqlx.Session) AttachModel {
	return NewAttachModel(sqlx.NewSqlConnFromSession(session))
}
