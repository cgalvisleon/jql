package jql

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jql/jql"
)

const (
	DriverPostgres = jql.DriverPostgres
	DriverSqlite   = jql.DriverSqlite
	// Field names
	SOURCE     = jql.SOURCE
	ID         = jql.ID
	IDX        = jql.IDX
	STATUS     = jql.STATUS
	VERSION    = jql.VERSION
	PROJECT_ID = jql.PROJECT_ID
	TENANT_ID  = jql.TENANT_ID
	CREATED_AT = jql.CREATED_AT
	UPDATED_AT = jql.UPDATED_AT
	// Data types
	ANY      = jql.ANY
	BYTES    = jql.BYTES
	INT      = jql.INT
	FLOAT    = jql.FLOAT
	KEY      = jql.KEY
	TEXT     = jql.TEXT
	MEMO     = jql.MEMO
	JSON     = jql.JSON
	DATETIME = jql.DATETIME
	BOOLEAN  = jql.BOOLEAN
	GEOMETRY = jql.GEOMETRY
	CALC     = jql.CALC
	// Column types
	COLUMN = jql.COLUMN
	ATTRIB = jql.ATTRIB
	DETAIL = jql.DETAIL
	ROLLUP = jql.ROLLUP
	AGG    = jql.AGG
	// Status record
	ACTIVE     = jql.ACTIVE
	ARCHIVED   = jql.ARCHIVED
	CANCELED   = jql.CANCELED
	OF_SYSTEM  = jql.OF_SYSTEM
	FOR_DELETE = jql.FOR_DELETE
	PENDING    = jql.PENDING
	APPROVED   = jql.APPROVED
	REJECTED   = jql.REJECTED
)

var (
	// Error
	ErrDuplicate   = jql.ErrDuplicate
	ErrNotInserted = jql.ErrNotInserted
	ErrNotUpdated  = jql.ErrNotUpdated
	ErrNotFound    = jql.ErrNotFound
	ErrNotUpserted = jql.ErrNotUpserted
)

type TypeColumn = jql.TypeColumn
type TypeData = jql.TypeData
type Driver = jql.Driver
type DB = jql.DB
type Model = jql.Model
type Tx = jql.Tx
type Condition = jql.Condition
type Ql = jql.Ql
type Cmd = jql.Cmd

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
* @param name string
* @return (*DB, error)
**/
func LoadTo(name string) (*DB, error) {
	return jql.LoadTo(name)
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
