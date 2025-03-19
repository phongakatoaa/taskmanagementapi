package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type (
	DB struct {
		pg *sql.DB
	}

	FindOptions struct {
		AssignedUserIDs []int
		Statuses        []Status
		SortBy          DBColumn
		SortOrder       SortOrder
	}

	DBColumn string

	SortOrder string
)

const (
	IDCol             DBColumn = "tasks.id"
	TitleCol          DBColumn = "tasks.title"
	DescriptionCol    DBColumn = "tasks.description"
	AssignedUserIDCol DBColumn = "tasks.assigned_user_id"
	StatusCol         DBColumn = "tasks.status"
	CreatedAtCol      DBColumn = "tasks.created_at"
	DueDateCol        DBColumn = "tasks.due_date"

	AssignedUsernameCol DBColumn = "users.username AS assigned_username"
)

var allColumns = []DBColumn{IDCol, TitleCol, DescriptionCol, AssignedUserIDCol, StatusCol, CreatedAtCol, DueDateCol}

func (col DBColumn) String() string {
	return string(col)
}

func ParseDBColumn(str string) (DBColumn, error) {
	if str == "assigned_username" {
		return DBColumn(str), nil
	}
	if !strings.HasPrefix(str, "tasks") {
		str = fmt.Sprintf("tasks.%s", str)

	}
	for _, col := range allColumns {
		if str == col.String() {
			return col, nil
		}
	}
	return "", fmt.Errorf("invalid column: %s", str)
}

const (
	SortOrderAscending  SortOrder = "ASC"
	SortOrderDescending SortOrder = "DESC"
)

func ParseSortOrder(str string) (SortOrder, error) {
	if strings.EqualFold(str, string(SortOrderAscending)) {
		return SortOrderAscending, nil
	}
	if strings.EqualFold(str, string(SortOrderDescending)) {
		return SortOrderDescending, nil
	}
	return "", fmt.Errorf("invalid sort order: %s", str)
}

func NewDB(pg *sql.DB) *DB {
	return &DB{pg}
}

func (db *DB) Find(ctx context.Context, options FindOptions) ([]Entry, error) {
	query, args := options.buildQuery()
	rows, err := db.pg.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return db.scanRows(rows)
}

func (db *DB) Insert(ctx context.Context, entry Entry) (id int, err error) {
	if err := entry.Validate(); err != nil {
		return 0, fmt.Errorf("invalid entry: %w", err)
	}
	var assignedUserID *int
	if entry.AssignedUserID > 0 {
		assignedUserID = &entry.AssignedUserID
	}
	status := StatusPending
	if entry.Status != "" {
		status = entry.Status
	}

	row := db.pg.QueryRowContext(ctx,
		"INSERT INTO api.tasks (title,description,assigned_user_id,status,due_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		entry.Title, entry.Description, assignedUserID, status, entry.DueDate,
	)
	err = row.Scan(&id)
	return id, err
}

func (db *DB) UpdateStatus(ctx context.Context, id int, assignedUserID int, status Status) error {
	row := db.pg.QueryRowContext(ctx,
		"UPDATE api.tasks SET status = $1 WHERE assigned_user_id = $2 AND id = $3 RETURNING id",
		status, id, assignedUserID,
	)
	var updatedID int
	return row.Scan(&updatedID)
}

type TaskSummary struct {
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	Assigned  int    `json:"assigned"`
	Completed int    `json:"completed"`
}

func (db *DB) Summarize(ctx context.Context) ([]TaskSummary, error) {
	stmt := fmt.Sprintf(`
		SELECT 
			users.id,
			users.username,
			COUNT(*) as assigned,
			COUNT(*) FILTER (WHERE status = 'COMPLETED') as completed
		FROM api.tasks
		JOIN auth.users ON users.id = tasks.assigned_user_id
		GROUP BY users.id ORDER BY users.id ASC
	`)

	rows, err := db.pg.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var summaries []TaskSummary

	for rows.Next() {
		s := TaskSummary{}
		if err = rows.Scan(&s.UserID, &s.Username, &s.Assigned, &s.Completed); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}

func (db *DB) scanRows(rows *sql.Rows) ([]Entry, error) {
	var entries []Entry
	for rows.Next() {
		var entry Entry

		err := rows.Scan(
			&entry.ID, &entry.Title, &entry.Description, &entry.AssignedUserID,
			&entry.Status, &entry.CreatedAt, &entry.DueDate, &entry.AssignedUsername,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (opt FindOptions) buildQuery() (query string, args []interface{}) {
	var clauses []string
	args = make([]interface{}, 0)

	if len(opt.AssignedUserIDs) > 0 {
		args = append(args, pq.Array(opt.AssignedUserIDs))
		clauses = append(clauses, fmt.Sprintf("assigned_user_id = ANY($%d)", len(args)))
	}
	if len(opt.Statuses) > 0 {
		args = append(args, pq.Array(opt.Statuses))
		clauses = append(clauses, fmt.Sprintf("status = ANY($%d)", len(args)))
	}

	selectCols := append(allColumns, AssignedUsernameCol)
	stmt := fmt.Sprintf(
		"SELECT %s FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id",
		strings.Join(toColumnStrings(selectCols), ","),
	)

	if len(clauses) == 0 {
		return fmt.Sprintf("%s%s", stmt, opt.sortClause()), args
	}
	return fmt.Sprintf("%s WHERE %s%s", stmt, strings.Join(clauses, " AND "), opt.sortClause()), args
}

func (opt FindOptions) sortClause() string {
	if opt.SortBy == "" {
		return ""
	}
	sortOrder := SortOrderAscending
	if opt.SortOrder == SortOrderDescending {
		sortOrder = SortOrderDescending
	}
	return fmt.Sprintf(" ORDER BY %s %s", opt.SortBy, sortOrder)
}

func toColumnStrings(cols []DBColumn) []string {
	var ret []string
	for _, col := range cols {
		ret = append(ret, col.String())
	}
	return ret
}
