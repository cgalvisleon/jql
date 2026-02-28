package jql

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jql/jdb"
	"github.com/cgalvisleon/jql/jql"
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
	ErrNotUpdated = jdb.ErrNotUpdated
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
* NewModel
* @param db *DB, schema, name string, version int
* @return (*Model, error)
**/
func NewModel(db *DB, schema, name string, version int) (*Model, error) {
	result, err := db.Define(et.Json{
		"schema":  schema,
		"name":    name,
		"version": version,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* ConnectTo
* @param name string, params et.Json
* @return (*DB, error)
**/
func ConnectTo(name string, params et.Json) (*DB, error) {
	return jql.ConnectTo(name, params)
}

/**
* LoadTo
* @param name string, host string, port int
* @return (*DB, error)
**/
func LoadTo(name, host string, port int) (*DB, error) {
	return jql.LoadTo(name, host, port)
}

/**
* Load
* @return (*DB, error)
**/
func Load() (*DB, error) {
	return jql.Load()
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func Define(definition et.Json) (*Model, error) {
	return jql.Define(definition)
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func Query(query et.Json) (et.Items, error) {
	return jql.Query(query)
}
