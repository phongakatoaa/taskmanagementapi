package auth

import (
	"github.com/lib/pq"
	"reflect"
	"testing"
)

func TestFindOptions_BuildQuery(t *testing.T) {
	cases := []struct {
		name   string
		option FindOptions
		query  string
		args   []interface{}
	}{
		{
			name:   "no options",
			option: FindOptions{},
			query:  "SELECT users.id,users.username,users.password,users.created_at,users.role FROM auth.users",
			args:   []interface{}{},
		},
		{
			name: "with IDs",
			option: FindOptions{
				IDs: []int{1, 2},
			},
			query: "SELECT users.id,users.username,users.password,users.created_at,users.role FROM auth.users WHERE id = ANY($1)",
			args: []interface{}{
				pq.Array([]int{1, 2}),
			},
		},
		{
			name: "with roles",
			option: FindOptions{
				Roles: []Role{RoleEmployee},
			},
			query: "SELECT users.id,users.username,users.password,users.created_at,users.role FROM auth.users WHERE role = ANY($1)",
			args: []interface{}{
				pq.Array([]Role{RoleEmployee}),
			},
		},
		{
			name: "with IDs and roles",
			option: FindOptions{
				IDs:   []int{1, 2},
				Roles: []Role{RoleEmployee},
			},
			query: "SELECT users.id,users.username,users.password,users.created_at,users.role FROM auth.users WHERE id = ANY($1) AND role = ANY($2)",
			args: []interface{}{
				pq.Array([]int{1, 2}),
				pq.Array([]Role{RoleEmployee}),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query, args := c.option.buildQuery()
			if query != c.query {
				t.Errorf("expected %s, got %s", c.query, query)
			}
			if len(args) != len(c.args) {
				t.Errorf("expected %d args, got %d", len(c.args), len(args))
			}
			for i, arg := range args {
				if !reflect.DeepEqual(arg, c.args[i]) {
					t.Errorf("expected %v, got %v", c.args[i], arg)
				}
			}
		})
	}
}
