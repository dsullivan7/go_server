package utilities

import (
	"gorm.io/gorm"
)

type DatabaseUtility interface {
	TruncateAll()
}

type GormDatabaseUtility struct {
	database *gorm.DB
}

func NewGormDatabaseUtility(database *gorm.DB) DatabaseUtility {
	return &GormDatabaseUtility{database: database}
}

func (gormDatabaseUtility *GormDatabaseUtility) TruncateAll() {
	gormDatabaseUtility.database.Exec(`
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
}
