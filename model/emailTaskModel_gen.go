// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	emailTaskFieldNames          = builder.RawFieldNames(&EmailTask{})
	emailTaskRows                = strings.Join(emailTaskFieldNames, ",")
	emailTaskRowsExpectAutoSet   = strings.Join(stringx.Remove(emailTaskFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	emailTaskRowsWithPlaceHolder = strings.Join(stringx.Remove(emailTaskFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	emailTaskModel interface {
		Insert(ctx context.Context, data *EmailTask) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*EmailTask, error)
		FindAll(ctx context.Context,email string) ([]*EmailTask, error)
		Update(ctx context.Context, data *EmailTask) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultEmailTaskModel struct {
		conn  sqlx.SqlConn
		table string
	}

	EmailTask struct {
		Id         uint64         `db:"id"`
		Email      string         `db:"email"`       // 邮件地址
		ContentId  uint64 `db:"content_id"`  // 邮件内容id
		SendTime   int64         `db:"send_time"`   // 发送时间
		CreateTime  string   `db:"create_time"` // 创建时间
		UpdateTime  string  `db:"update_time"` // 更新时间
	}
)

func newEmailTaskModel(conn sqlx.SqlConn) *defaultEmailTaskModel {
	return &defaultEmailTaskModel{
		conn:  conn,
		table: "`email_task`",
	}
}

func (m *defaultEmailTaskModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultEmailTaskModel) FindOne(ctx context.Context, id uint64) (*EmailTask, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", emailTaskRows, m.table)
	var resp EmailTask
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
func (m *defaultEmailTaskModel) FindAll(ctx context.Context,email string) ([]*EmailTask, error) {
	selectBuilder := sq.Select("*").From(m.tableName())

	if len(strings.Trim(email,"")) >0 {
		selectBuilder=selectBuilder.Where(sq.Like{"email":email})
	}

	query, args, err := selectBuilder.Limit(1000).ToSql()
	var resp []*EmailTask
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

func (m *defaultEmailTaskModel) Insert(ctx context.Context, data *EmailTask) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, emailTaskRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Email, data.ContentId, data.SendTime)
	return ret, err
}

func (m *defaultEmailTaskModel) Update(ctx context.Context, data *EmailTask) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, emailTaskRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Email, data.ContentId, data.SendTime, data.Id)
	return err
}

func (m *defaultEmailTaskModel) tableName() string {
	return m.table
}
