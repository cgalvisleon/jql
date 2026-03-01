package jql

import (
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jql/jdb"
)

const (
	DriverPostgres = jdb.DriverPostgres
	DriverSqlite   = jdb.DriverSqlite
	// Field names
	SOURCE     = jdb.SOURCE
	ID         = jdb.ID
	IDX        = jdb.IDX
	STATUS     = jdb.STATUS
	VERSION    = jdb.VERSION
	PROJECT_ID = jdb.PROJECT_ID
	TENANT_ID  = jdb.TENANT_ID
	CREATED_AT = jdb.CREATED_AT
	UPDATED_AT = jdb.UPDATED_AT
	// Data types
	ANY      = jdb.ANY
	BYTES    = jdb.BYTES
	INT      = jdb.INT
	FLOAT    = jdb.FLOAT
	KEY      = jdb.KEY
	TEXT     = jdb.TEXT
	MEMO     = jdb.MEMO
	JSON     = jdb.JSON
	DATETIME = jdb.DATETIME
	BOOLEAN  = jdb.BOOLEAN
	GEOMETRY = jdb.GEOMETRY
	CALC     = jdb.CALC
	// Column types
	COLUMN = jdb.COLUMN
	ATTRIB = jdb.ATTRIB
	DETAIL = jdb.DETAIL
	ROLLUP = jdb.ROLLUP
	AGG    = jdb.AGG
	// Status record
	ACTIVE     = jdb.ACTIVE
	ARCHIVED   = jdb.ARCHIVED
	CANCELED   = jdb.CANCELED
	OF_SYSTEM  = jdb.OF_SYSTEM
	FOR_DELETE = jdb.FOR_DELETE
	PENDING    = jdb.PENDING
	APPROVED   = jdb.APPROVED
	REJECTED   = jdb.REJECTED
)

var (
	// Error
	ErrNotUpdated        = jdb.ErrNotUpdated
	ErrNotInserted error = fmt.Errorf("record not inserted")
	ErrNotFound    error = fmt.Errorf("record not found")
	ErrNotUpserted error = fmt.Errorf("record not inserted or updated")
	ErrDuplicate   error = fmt.Errorf("record duplicate")
)

type TypeColumn = jdb.TypeColumn
type TypeData = jdb.TypeData
type Driver = jdb.Driver
type DB = jdb.DB
type Model = jdb.Model
type Tx = jdb.Tx
type Condition = jdb.Condition
type Ql = jdb.Ql
type Cmd = jdb.Cmd

/**
* ConnectTo
* @param name string, params Connection
* @return *jdb.DB, error
**/
func ConnectTo(name string, params et.Json) (*jdb.DB, error) {
	return jdb.LoadDb(name, params)
}

/**
* LoadTo
* @param name string
* @return *jdb.DB, error
**/
func LoadTo(name string, host string, port int) (*jdb.DB, error) {
	params := et.Json{
		"driver":   envar.GetStr("DB_DRIVER", "postgres"),
		"database": name,
		"host":     host,
		"port":     port,
		"username": envar.GetStr("DB_USERNAME", "test"),
		"password": envar.GetStr("DB_PASSWORD", "test"),
		"app":      envar.GetStr("DB_APP", "jql"),
		"version":  envar.GetInt("DB_VERSION", 15),
	}

	return ConnectTo(name, params)
}

/**
* Load
* @return (*jdb.DB, error)
**/
func Load() (*jdb.DB, error) {
	name := envar.GetStr("DB_NAME", "josephine")
	host := envar.GetStr("DB_HOST", "localhost")
	port := envar.GetInt("DB_PORT", 5432)
	return LoadTo(name, host, port)
}

/**
* Define
* @param params et.Json
* @return *jdb.Model, error
**/
func Define(params et.Json) (*jdb.Model, error) {
	database := params.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(jdb.MSG_ATTRIBUTE_REQUIRED, "database")
	}

	db, err := jdb.GetDb(database)
	if err != nil {
		return nil, err
	}

	return db.Define(params)
}

/**
* Query
* @param params et.Json
* @return (et.Items, error)
**/
func Query(params et.Json) (et.Items, error) {
	return et.Items{}, nil
}
