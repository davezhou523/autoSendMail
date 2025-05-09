// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.6

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	searchConfigFieldNames          = builder.RawFieldNames(&SearchConfig{})
	searchConfigRows                = strings.Join(searchConfigFieldNames, ",")
	searchConfigRowsExpectAutoSet   = strings.Join(stringx.Remove(searchConfigFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	searchConfigRowsWithPlaceHolder = strings.Join(stringx.Remove(searchConfigFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	searchConfigModel interface {
		Insert(ctx context.Context, data *SearchConfig) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*SearchConfig, error)
		Update(ctx context.Context, data *SearchConfig) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultSearchConfigModel struct {
		conn  sqlx.SqlConn
		table string
	}

	SearchConfig struct {
		Id          uint64 `db:"id"`
		Key         string `db:"key"`
		EngineId    string `db:"engine_id"`
		ProjectName string `db:"project_name"`
		ProjectId   string `db:"project_id"`
		Sort        uint64 `db:"sort"` // 排序
	}
)

func newSearchConfigModel(conn sqlx.SqlConn) *defaultSearchConfigModel {
	return &defaultSearchConfigModel{
		conn:  conn,
		table: "`search_config`",
	}
}

func (m *defaultSearchConfigModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultSearchConfigModel) FindOne(ctx context.Context, id uint64) (*SearchConfig, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", searchConfigRows, m.table)
	var resp SearchConfig
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSearchConfigModel) Insert(ctx context.Context, data *SearchConfig) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, searchConfigRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Key, data.EngineId, data.ProjectName, data.ProjectId, data.Sort)
	return ret, err
}

func (m *defaultSearchConfigModel) Update(ctx context.Context, data *SearchConfig) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, searchConfigRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Key, data.EngineId, data.ProjectName, data.ProjectId, data.Sort, data.Id)
	return err
}

func (m *defaultSearchConfigModel) tableName() string {
	return m.table
}
