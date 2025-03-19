package tasks

import (
	"context"
	"reflect"
	"testing"

	"github.com/lib/pq"

	"siransbach/taskmanagementapi/postgres"
)

func TestFindOptions_BuildQuery(t *testing.T) {
	cases := []struct {
		name  string
		opts  FindOptions
		query string
		args  []interface{}
	}{
		{
			name: "no options",
			opts: FindOptions{},
			query: "SELECT tasks.id,tasks.title,tasks.description,tasks.assigned_user_id,tasks.status" +
				",tasks.created_at,tasks.due_date,users.username AS assigned_username" +
				" FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id",
			args: []interface{}{},
		},
		{
			name: "with assigned user IDs",
			opts: FindOptions{
				AssignedUserIDs: []int{1, 2},
			},
			query: "SELECT tasks.id,tasks.title,tasks.description,tasks.assigned_user_id,tasks.status" +
				",tasks.created_at,tasks.due_date,users.username AS assigned_username" +
				" FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id WHERE assigned_user_id = ANY($1)",
			args: []interface{}{
				pq.Array([]int{1, 2}),
			},
		}, {
			name: "with statuses",
			opts: FindOptions{
				Statuses: []Status{StatusCompleted, StatusInProgress},
			},
			query: "SELECT tasks.id,tasks.title,tasks.description,tasks.assigned_user_id,tasks.status" +
				",tasks.created_at,tasks.due_date,users.username AS assigned_username" +
				" FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id WHERE status = ANY($1)",
			args: []interface{}{
				pq.Array([]Status{StatusCompleted, StatusInProgress}),
			},
		},
		{
			name: "with assigned user IDs, statuses",
			opts: FindOptions{
				AssignedUserIDs: []int{1, 2},
				Statuses:        []Status{StatusCompleted, StatusInProgress},
			},
			query: "SELECT tasks.id,tasks.title,tasks.description,tasks.assigned_user_id,tasks.status" +
				",tasks.created_at,tasks.due_date,users.username AS assigned_username" +
				" FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id WHERE assigned_user_id = ANY($1) AND status = ANY($2)",
			args: []interface{}{
				pq.Array([]int{1, 2}),
				pq.Array([]Status{StatusCompleted, StatusInProgress}),
			},
		},
		{
			name: "with sort",
			opts: FindOptions{
				SortBy:    CreatedAtCol,
				SortOrder: SortOrderAscending,
			},
			query: "SELECT tasks.id,tasks.title,tasks.description,tasks.assigned_user_id,tasks.status" +
				",tasks.created_at,tasks.due_date,users.username AS assigned_username" +
				" FROM api.tasks JOIN auth.users ON users.id = tasks.assigned_user_id ORDER BY tasks.created_at ASC",
			args: []interface{}{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query, args := tc.opts.buildQuery()
			if query != tc.query {
				t.Errorf("expected query %q, got %q", tc.query, query)
			}
			if !reflect.DeepEqual(args, tc.args) {
				t.Errorf("expected args %v, got %v", tc.args, args)
			}
		})
	}
}

func TestDB_Find(t *testing.T) {
	pg, err := postgres.Connect(postgres.ConnectionString())
	if err != nil {
		panic(err)
	}

	entries, err := NewDB(pg).Find(context.Background(), FindOptions{
		AssignedUserIDs: []int{1},
	})
	if err != nil {
		panic(err)
	}
	t.Log(entries)
}

func TestDB_UpdateStatus(t *testing.T) {
	pg, err := postgres.Connect(postgres.ConnectionString())
	if err != nil {
		panic(err)
	}

	err = NewDB(pg).UpdateStatus(context.Background(), 1, 1, StatusCompleted)
	if err != nil {
		panic(err)
	}
}
