package utilities

import (
	"database/sql"
)

type DatabaseUtility interface {
	TruncateAll() error
}

type SQLDatabaseUtility struct {
	database *sql.DB
}

func NewSQLDatabaseUtility(database *sql.DB) DatabaseUtility {
	return &SQLDatabaseUtility{database: database}
}

func (dbUtility *SQLDatabaseUtility) TruncateAll() error {
	_, err := dbUtility.database.Exec(`
		do $$
		begin
			execute (
				select 'truncate table ' || string_agg('"' || tablename || '"', ', ')
				from pg_tables
				where schemaname = 'public'
			);
		end;
		$$
	`)

  return err
}
