package model

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ AttachModel = (*customAttachModel)(nil)

type (
	// AttachModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAttachModel.
	AttachModel interface {
		attachModel
		withSession(session sqlx.Session) AttachModel
		FindAll(ctx context.Context, idjson string) ([]*Attach, error)
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

func (m *customAttachModel) FindAll(ctx context.Context, idjson string) ([]*Attach, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if len(strings.Trim(idjson, "")) > 0 {
		var ids []int
		err := json.Unmarshal([]byte(idjson), &ids)
		if err != nil {
			fmt.Printf("attach_id json失败: %v\n", err)
			return nil, err
		}
		selectBuilder = selectBuilder.Where(sq.Eq{"id": ids})
	}

	query, args, err := selectBuilder.Limit(1000).ToSql()
	var resp []*Attach
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
