package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateUsers_20180728_230541 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateUsers_20180728_230541{}
	m.Created = "20180728_230541"

	migration.Register("CreateUsers_20180728_230541", m)
}

// Up Run the migrations
func (m *CreateUsers_20180728_230541) Up() {
	sql := `
		CREATE TABLE users(
			id serial primary key,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
	  )
	`
	m.SQL(sql)
}

// Down Reverse the migrations
func (m *CreateUsers_20180728_230541) Down() {
	m.SQL("DROP TABLE users")
}
