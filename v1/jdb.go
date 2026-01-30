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
	KEY        = jql.KEY
	IDX        = jql.IDX
	STATUS     = jql.STATUS
	VERSION    = jql.VERSION
	PROJECT_ID = jql.PROJECT_ID
	TENANT_ID  = jql.TENANT_ID
	CREATED_AT = jql.CREATED_AT
	UPDATED_AT = jql.UPDATED_AT
	// Data types
	TpAny      = jql.TpAny
	TpBytes    = jql.TpBytes
	TpInt      = jql.TpInt
	TpFloat    = jql.TpFloat
	TpKey      = jql.TpKey
	TpText     = jql.TpText
	TpMemo     = jql.TpMemo
	TpJson     = jql.TpJson
	TpDateTime = jql.TpDateTime
	TpBoolean  = jql.TpBoolean
	TpGeometry = jql.TpGeometry
	TpCalc     = jql.TpCalc
	// Column types
	TpColumn      = jql.TpColumn
	TpAtrib       = jql.TpAtrib
	TpDetail      = jql.TpDetail
	TpRollup      = jql.TpRollup
	TpAggregation = jql.TpAggregation
	// Status record
	Active    = jql.Active
	Archived  = jql.Archived
	Canceled  = jql.Canceled
	OfSystem  = jql.OfSystem
	ForDelete = jql.ForDelete
	Pending   = jql.Pending
	Approved  = jql.Approved
	Rejected  = jql.Rejected
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
* @param tenantId, name, driver string, userCore bool, params et.Json
* @return (*DB, error)
**/
func ConnectTo(tenantId, name, driver string, userCore bool, params et.Json) (*DB, error) {
	return jql.ConnectTo(tenantId, name, driver, userCore, params)
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
* From
* @param model *Model, as string
* @return *Ql
**/
func From(model *Model, as string) *Ql {
	return jql.From(model, as)
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func Query(query et.Json) (et.Items, error) {
	return jql.Query(query)
}
