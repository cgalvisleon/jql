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

type DriverFn func() Driver

var drivers map[string]DriverFn

func init() {
	drivers = make(map[string]DriverFn)
}

func Register(name string, driver DriverFn) {
	drivers[name] = driver
}
