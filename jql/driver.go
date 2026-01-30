package jql

import "database/sql"

const (
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite"
)

type Driver interface {
	Connect(db *DB) (*sql.DB, error)
	Load(model *Model) (string, error)
	Mutate(model *Model) (string, error)
	Query(query *Ql) (string, error)
	Command(command *Cmd) (string, error)
}

var drivers map[string]Driver

func init() {
	drivers = make(map[string]Driver)
}

func Register(name string, driver Driver) {
	drivers[name] = driver
}
