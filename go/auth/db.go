package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type (
	DB struct {
		session *sql.DB
	}

	DBColumn string

	FindOptions struct {
		IDs      []int
		Username string
		Roles    []Role
	}
)

const (
	IDCol        DBColumn = "users.id"
	UsernameCol  DBColumn = "users.username"
	PasswordCol  DBColumn = "users.password"
	CreatedAtCol DBColumn = "users.created_at"
	RoleCol      DBColumn = "users.role"
)

var allCols = []DBColumn{IDCol, UsernameCol, PasswordCol, CreatedAtCol, RoleCol}

func NewDB(session *sql.DB) *DB {
	return &DB{session}
}

func (db *DB) Find(ctx context.Context, options FindOptions) ([]*User, error) {
	query, args := options.buildQuery()
	rows, err := db.session.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return db.scanRows(rows)
}

func (db *DB) FindOne(ctx context.Context, options FindOptions) (*User, error) {
	query, args := options.buildQuery()

	var user User
	err := db.session.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Username, &user.EncryptedPassword, &user.CreatedAt, &user.Role,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) scanRows(rows *sql.Rows) ([]*User, error) {
	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.EncryptedPassword, &user.CreatedAt, &user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (col DBColumn) String() string {
	return string(col)
}

func (opt FindOptions) buildQuery() (query string, args []interface{}) {
	var whereConds []string
	if len(opt.IDs) > 0 {
		args = append(args, pq.Array(opt.IDs))
		whereConds = append(whereConds, fmt.Sprintf("id = ANY($%d)", len(args)))
	}
	if len(opt.Username) > 0 {
		args = append(args, opt.Username)
		whereConds = append(whereConds, fmt.Sprintf("username = $%d", len(args)))
	}
	if len(opt.Roles) > 0 {
		args = append(args, pq.Array(opt.Roles))
		whereConds = append(whereConds, fmt.Sprintf("role = ANY($%d)", len(args)))
	}
	if len(whereConds) == 0 {
		return fmt.Sprintf(
			"SELECT %s FROM auth.users",
			strings.Join(toColumnStrings(allCols), ","),
		), args
	}
	return fmt.Sprintf(
		"SELECT %s FROM auth.users WHERE %s",
		strings.Join(toColumnStrings(allCols), ","),
		strings.Join(whereConds, " AND "),
	), args
}

func toColumnStrings(cols []DBColumn) []string {
	var ret []string
	for _, col := range cols {
		ret = append(ret, col.String())
	}
	return ret
}
