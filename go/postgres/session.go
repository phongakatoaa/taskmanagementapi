package postgres

import "database/sql"

func Connect(connectionString string) (*sql.DB, error) {
	return sql.Open("postgres", connectionString)
}
