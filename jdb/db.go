package jdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Schema struct {
	Database string           `json:"-"`
	Name     string           `json:"name"`
	Models   map[string]*From `json:"models"`
}

type DB struct {
	Name    string             `json:"name"`
	Schemas map[string]*Schema `json:"schemas"`
	Params  et.Json            `json:"params"`
	driver  Driver             `json:"-"`
	db      *sql.DB            `json:"-"`
}

/**
* serialize
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
* Save
* @return error
**/
func (s *DB) Save() error {
	bt, err := s.serialize()
	if err != nil {
		return err
	}

	return setCatalog("db", s.Name, 1, bt)
}

/**
* init
* @return error
**/
func (s *DB) init() error {
	if s.driver == nil {
		return errors.New(MSG_DRIVER_NOT_FOUND)
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
	return s.Save()
}

/**
* NewModel
* @param schema, name string, version int
* @return *Model
**/
func (s *DB) NewModel(schema, name string, version int) (*Model, error) {
	key := name
	key = strs.Append(schema, key, ".")
	key = strs.Append(s.Name, key, ".")

	result, ok := models[key]
	if ok {
		return result, nil
	}

	schema = utility.Normalize(schema)
	name = utility.Normalize(name)
	result = &Model{
		Database:      s.Name,
		Schema:        schema,
		Name:          name,
		Columns:       make([]*Column, 0),
		Indexes:       make([]string, 0),
		PrimaryKeys:   make([]string, 0),
		ForeignKeys:   make([]*Detail, 0),
		Unique:        make([]string, 0),
		Required:      make([]string, 0),
		Hidden:        make([]string, 0),
		Details:       make(map[string]*Detail, 0),
		Rollups:       make(map[string]*Detail, 0),
		Relations:     make(map[string]*Detail, 0),
		Version:       version,
		beforeInserts: make([]TriggerFunction, 0),
		beforeUpdates: make([]TriggerFunction, 0),
		beforeDeletes: make([]TriggerFunction, 0),
		afterInserts:  make([]TriggerFunction, 0),
		afterUpdates:  make([]TriggerFunction, 0),
		afterDeletes:  make([]TriggerFunction, 0),
		calcs:         make(map[string]DataContext),
	}
	result.defineIdxField()

	sch := s.getSchema(schema)
	sch.Models[name] = result.From()
	models[key] = result

	return result, nil
}

/**
* sqlTx
* @param tx *Tx, sql string, arg ...any
* @return et.Items, error
*
 */
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
		result := RowsToItems(rows)
		return result, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	result := RowsToItems(rows)
	return result, nil
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
		Models:   make(map[string]*From),
	}
	s.Schemas[name] = result
	return result
}

/**
* GetModel
* @param schema, name string
* @return *Model
**/
func (s *DB) GetModel(schema, name string) (*Model, error) {
	key := name
	key = strs.Append(schema, key, ".")
	key = strs.Append(s.Name, key, ".")

	result, ok := models[key]
	if ok {
		return result, nil
	}

	exists, err := getCatalog("model", key, &result)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrModelNotFound
	}

	result.beforeInserts = make([]TriggerFunction, 0)
	result.beforeUpdates = make([]TriggerFunction, 0)
	result.beforeDeletes = make([]TriggerFunction, 0)
	result.afterInserts = make([]TriggerFunction, 0)
	result.afterUpdates = make([]TriggerFunction, 0)
	result.afterDeletes = make([]TriggerFunction, 0)
	result.calcs = make(map[string]DataContext)
	result.db = s

	err = result.Init()
	if err != nil {
		return nil, err
	}

	sch := s.getSchema(schema)
	sch.Models[name] = result.From()
	models[key] = result

	return result, nil
}

/**
* DeleteModel
* @param schema, name string
* @return error
**/
func (s *DB) DeleteModel(schema, name string) error {
	key := name
	key = strs.Append(schema, key, ".")
	key = strs.Append(s.Name, key, ".")

	_, ok := models[key]
	if ok {
		delete(models, key)
	}

	err := deleteCatalog("model", key)
	if err != nil {
		return err
	}

	return nil
}

/**
* Command
* @param command *Command
* @return et.Items, error
**/
func (s *DB) Command(command *Cmd) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, errors.New(MSG_DRIVER_NOT_FOUND)
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
		return et.Items{}, errors.New(MSG_DRIVER_NOT_FOUND)
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
func (s *DB) From(sql et.Json) (et.Items, error) {
	return et.Items{}, nil
}
