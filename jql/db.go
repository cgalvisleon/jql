package jql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/utility"
)

var dbs map[string]*DB

func init() {
	dbs = make(map[string]*DB)
}

type DB struct {
	Name    string             `json:"name"`
	Schemas map[string]*Schema `json:"schemas"`
	Params  et.Json            `json:"params"`
	driver  Driver             `json:"-"`
	db      *sql.DB            `json:"-"`
}

/**
* getDb
* @param name string, params Connection
* @return (*DB, error)
**/
func getDb(name string, params et.Json) (*DB, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "name")
	}

	name = utility.Normalize(name)
	result, ok := dbs[name]
	if ok {
		return result, nil
	}

	driver := params.Str("driver")
	drv, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
	}

	result = &DB{
		Name:    name,
		Schemas: make(map[string]*Schema),
		Params:  params,
		driver:  drv(),
	}
	err := result.load()
	if err != nil {
		return nil, err
	}

	dbs[name] = result
	return result, nil
}

/**
* Serialize
* @return []byte, error
**/
func (s *DB) serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *DB) ToJson() et.Json {
	bt, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* load
* @return error
**/
func (s *DB) load() error {
	if s.driver == nil {
		return fmt.Errorf(MSG_DRIVER_NOT_FOUND)
	}

	db, err := s.driver.Connect(s)
	if err != nil {
		return err
	}

	isCore := s.Params.Bool("is_core")
	if !isCore {
		s.initCore()
	}

	s.db = db
	return nil
}

/**
* getSchema
* @param name string
* @return *Schema
**/
func (s *DB) getSchema(name string) *Schema {
	result, ok := s.Schemas[name]
	if ok {
		return result
	}

	result = &Schema{
		Database: s.Name,
		Name:     name,
		Models:   make(map[string]*Model),
		Db:       s,
	}
	s.Schemas[name] = result
	return result
}

/**
* NewModel
* @param schema, name string, version int
* @return *Model
**/
func (s *DB) NewModel(schema, name string, version int) (*Model, error) {
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, schema)
	}

	sch := s.getSchema(schema)
	return sch.newModel(name, version)
}

/**
* GetModel
* @param schema, name string
* @return *Model
**/
func (s *DB) GetModel(schema, name string) (*Model, error) {
	sch, ok := s.Schemas[schema]
	if !ok {
		return nil, fmt.Errorf(MSG_SCHEMA_NOT_FOUND, schema)
	}

	return sch.getModel(name)
}

/**
* deleteModel
* @param name string
* @return error
**/
func (s *DB) deleteModel(schema, name string) error {
	sch, ok := s.Schemas[schema]
	if !ok {
		return fmt.Errorf(MSG_SCHEMA_NOT_FOUND, schema)
	}

	return sch.deleteModel(name)
}

/**
* sqlTx
* @param tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) sqlTx(tx *Tx, _sql string, arg ...any) (et.Items, error) {
	query := SQLParse(_sql, arg...)
	if tx != nil {
		err := tx.Begin(s.db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err := tx.Tx.Query(query)
		if err != nil {
			errR := tx.Rollback()
			if errR != nil {
				err = fmt.Errorf(MSG_ROLLBACK_ERROR, errR)
			}
			return et.Items{}, err
		}
		result := rowsToItems(rows)
		return result, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	result := rowsToItems(rows)
	return result, nil
}

/**
* Command
* @param command *Command
* @return et.Items, error
**/
func (s *DB) Command(command *Cmd) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, fmt.Errorf(MSG_DRIVER_NOT_FOUND)
	}

	if command.IsDebug {
		logs.Debugf("command:%s", command.ToJson().ToEscapeHTML())
	}

	sql, err := s.driver.Command(command)
	if err != nil {
		return et.Items{}, err
	}

	return s.sqlTx(command.tx, sql)
}

/**
* Ql
* @param query *Ql
* @return et.Items, error
**/
func (s *DB) Query(query *Ql) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, fmt.Errorf(MSG_DRIVER_NOT_FOUND)
	}

	if query.IsDebug {
		logs.Debugf("query:%s", query.ToJson().ToEscapeHTML())
	}

	sql, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	result, err := s.sqlTx(query.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	wg := &sync.WaitGroup{}
	for _, item := range result.Result {
		wg.Add(1)
		go func(item et.Json) {
			query.getDetails(query.tx, item)
			query.getRollups(query.tx, item)
			query.getCalls(query.tx, item)
		}(item)
	}
	wg.Wait()

	return result, nil
}

/**
* Define
* @return *Model, error
**/
func (s *DB) Define(definition et.Json) (*Model, error) {
	return nil, nil
}

/**
* Insert
* @param insert et.Json
* @return et.Items, error
**/
func (s *DB) Insert(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* Update
* @param insert et.Json
* @return et.Items, error
**/
func (s *DB) Update(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* Delete
* @param insert et.Json
* @return et.Items, error
**/
func (s *DB) Delete(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* Delete
* @param insert et.Json
* @return et.Items, error
**/
func (s *DB) Upsert(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* Delete
* @param insert et.Json
* @return et.Items, error
**/
func (s *DB) Select(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}
