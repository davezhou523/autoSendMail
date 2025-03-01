package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CompanyModel = (*customCompanyModel)(nil)

type (
	// CompanyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCompanyModel.
	CompanyModel interface {
		companyModel
		withSession(session sqlx.Session) CompanyModel
	}

	customCompanyModel struct {
		*defaultCompanyModel
	}
)

// NewCompanyModel returns a model for the database table.
func NewCompanyModel(conn sqlx.SqlConn) CompanyModel {
	return &customCompanyModel{
		defaultCompanyModel: newCompanyModel(conn),
	}
}

func (m *customCompanyModel) withSession(session sqlx.Session) CompanyModel {
	return NewCompanyModel(sqlx.NewSqlConnFromSession(session))
}
