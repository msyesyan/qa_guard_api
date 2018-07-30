package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateProjectUsers_20180730_143558 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateProjectUsers_20180730_143558{}
	m.Created = "20180730_143558"

	migration.Register("CreateProjectUsers_20180730_143558", m)
}

// Run the migrations
func (m *CreateProjectUsers_20180730_143558) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
		CREATE TABLE project_users(
			id serial primary key,
			project_id integer NOT NULL,
			user_id integer NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)

	m.SQL("CREATE UNIQUE INDEX index_project_users_on_project_id on project_users USING btree (project_id)")
	m.SQL("CREATE UNIQUE INDEX index_project_users_on_user_id on project_users USING btree (user_id)")
}

// Reverse the migrations
func (m *CreateProjectUsers_20180730_143558) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
  m.SQL("DROP INDEX IF EXISTS index_project_users_on_project_id")
  m.SQL("DROP INDEX IF EXISTS index_project_users_on_user_id")
	m.SQL("DROP TABLE project_users")
}
