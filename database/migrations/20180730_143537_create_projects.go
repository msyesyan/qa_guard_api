package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateProjects_20180730_143537 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateProjects_20180730_143537{}
	m.Created = "20180730_143537"

	migration.Register("CreateProjects_20180730_143537", m)
}

// Run the migrations
func (m *CreateProjects_20180730_143537) Up() {
	sql := `
		CREATE TABLE projects(
			id serial primary key,
			title TEXT NOT NULL,
			code TEXT NOT NULL,
			user_id integer NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`
	m.SQL(sql)

	m.SQL("CREATE UNIQUE INDEX index_projects_on_user_id on projects USING btree (user_id)")
}

// Reverse the migrations
func (m *CreateProjects_20180730_143537) Down() {
	m.SQL("DROP INDEX IF EXISTS index_projects_on_user_id")
	m.SQL("DROP TABLE projects")
}
