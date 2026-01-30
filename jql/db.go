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
	Name       string             `json:"name"`
	Schemas    map[string]*Schema `json:"schemas"`
	Connection et.Json            `json:"connection"`
	driver     Driver             `json:"-"`
	db         *sql.DB            `json:"-"`
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
		Name:       name,
		Schemas:    make(map[string]*Schema),
		Connection: params,
		driver:     drv,
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
		db:       s,
	}
	s.Schemas[name] = result
	return result
}

/**
* newModel
* @param schema, name string, version int
* @return *Model
**/
func (s *DB) newModel(schema, name string, version int) (*Model, error) {
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, schema)
	}

	sch := s.getSchema(schema)
	return sch.newModel(name, version)
}

/**
* getModel
* @param schema, name string
* @return *Model
**/
func (s *DB) getModel(schema, name string) (*Model, error) {
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

	return s.SqlTx(command.tx, sql)
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
		logs.Debugf("command:%s", query.ToJson().ToEscapeHTML())
	}

	sql, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	result, err := s.SqlTx(query.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	wg := &sync.WaitGroup{}
	for _, item := range result.Result {
		wg.Add(1)
		go func(item et.Json) {
			query.getDetailsTx(query.tx, item)
			query.getRollupsTx(query.tx, item)
			query.getCallsTx(query.tx, item)
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
	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "schema")
	}

	name := definition.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "name")
	}

	version := definition.ValInt(1, "version")
	result, err := s.newModel(schema, name, version)
	if err != nil {
		return nil, err
	}

	isCore := definition.Bool("is_core")
	if isCore {
		result.IsCore = true
	}

	columns := definition.ArrayJson("columns")
	for _, column := range columns {
		name := column.String("name")
		tpColumn := column.String("tp_column")
		tpData := column.String("tp_data")
		hidden := column.Bool("hidden")
		defaultValue := column.String("default_value")
		definition, err := column.Byte("definition")
		if err != nil {
			return nil, err
		}
		_, err = result.defineColumn(name, TypeColumn(tpColumn), TypeData(tpData), hidden, defaultValue, definition)
		if err != nil {
			return nil, err
		}
	}

	source := definition.String("source")
	if utility.ValidStr(source, 0, []string{}) {
		err := result.DefineSourceField(source)
		if err != nil {
			return nil, err
		}
	}

	indexField := definition.String("index_field")
	if utility.ValidStr(indexField, 0, []string{}) {
		err := result.DefineIndexField(indexField)
		if err != nil {
			return nil, err
		}
	}

	primaryKeys := definition.ArrayStr("primary_keys")
	for _, primaryKey := range primaryKeys {
		result.DefinePrimaryKeys(primaryKey)
	}

	unique := definition.ArrayStr("unique")
	for _, unique := range unique {
		result.DefineUnique(unique)
	}

	indexes := definition.ArrayStr("indexes")
	for _, index := range indexes {
		err := result.DefineIndex(index)
		if err != nil {
			return nil, err
		}
	}

	required := definition.ArrayStr("required")
	for _, required := range required {
		result.DefineRequired(required)
	}

	hidden := definition.ArrayStr("hidden")
	for _, hidden := range hidden {
		result.DefineHidden(hidden)
	}

	details := definition.Json("details")
	for name := range details {
		jkeys := details.Json(name)
		keys := make(map[string]string, 0)
		version := details.ValInt(1, "version")
		for pk := range jkeys {
			fk := jkeys.Str(pk)
			keys[pk] = fk
		}
		_, err = result.DefineDetail(name, keys, version)
		if err != nil {
			return nil, err
		}
	}

	rollups := definition.Json("rollups")
	for name := range rollups {
		detail := rollups.Json(name)
		from := detail.String("from", "name")
		jkeys := detail.Json("keys")
		keys := make(map[string]string, 0)
		for pk := range jkeys {
			fk := jkeys.Str(pk)
			keys[pk] = fk
		}
		selects := detail.ArrayStr("selects")
		err = result.DefineRollup(name, from, keys, selects)
		if err != nil {
			return nil, err
		}
	}

	relations := definition.Json("relations")
	for name := range relations {
		detail := rollups.Json(name)
		jkeys := detail.Json("keys")
		keys := make(map[string]string, 0)
		for pk := range jkeys {
			fk := jkeys.Str(pk)
			keys[pk] = fk
		}
		err = result.DefineRelation(name, keys)
		if err != nil {
			return nil, err
		}
	}

	result.IsStrict = definition.Bool("is_strict")
	result.IsDebug = definition.Bool("is_debug")
	build := definition.Bool("build")
	if !build {
		return result, nil
	}

	if err = result.Init(); err != nil {
		return nil, err
	}

	return result, nil
}

/**
* Select
* @param query et.Json
* @return et.Items, error
**/
func (s *DB) Select(query et.Json) (et.Items, error) {
	bt, err := query.ToByte()
	if err != nil {
		return et.Items{}, err
	}

	var ql *Ql
	err = json.Unmarshal(bt, &ql)
	if err != nil {
		return et.Items{}, err
	}

	ql.DB = s
	result, err := ql.All()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Insert
* @param query et.Json
* @return et.Items, error
**/
func (s *DB) Insert(query et.Json) (et.Items, error) {
	from := query.String("from")
	if !utility.ValidStr(from, 0, []string{}) {
		return et.Items{}, fmt.Errorf(MSG_FROM_REQUIRED)
	}

	model, err := s.getModel(schema, from)
	if err != nil {
		return et.Items{}, err
	}

	data := query.Json("data")
	if len(data) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	cmd := model.Insert(data)
	result, err := cmd.Exec()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Update
* @param query et.Json
* @return et.Items, error
**/
func (s *DB) Update(query et.Json) (et.Items, error) {
	from := query.String("from")
	if !utility.ValidStr(from, 0, []string{}) {
		return et.Items{}, fmt.Errorf(MSG_FROM_REQUIRED)
	}

	model, err := s.GetModel(from)
	if err != nil {
		return et.Items{}, err
	}

	data := query.Json("data")
	if len(data) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	where := query.ArrayJson("where")
	if len(where) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	cmd := model.Update(data)
	cmd.WhereByJson(where)
	result, err := cmd.Exec()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Delete
* @param query et.Json
* @return et.Items, error
**/
func (s *DB) Delete(query et.Json) (et.Items, error) {
	from := query.String("from")
	if !utility.ValidStr(from, 0, []string{}) {
		return et.Items{}, fmt.Errorf(MSG_FROM_REQUIRED)
	}

	model, err := s.GetModel(from)
	if err != nil {
		return et.Items{}, err
	}

	where := query.ArrayJson("where")
	if len(where) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	cmd := model.Delete()
	cmd.WhereByJson(where)
	result, err := cmd.Exec()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Upsert
* @param query et.Json
* @return et.Items, error
**/
func (s *DB) Upsert(query et.Json) (et.Items, error) {
	from := query.String("from")
	if !utility.ValidStr(from, 0, []string{}) {
		return et.Items{}, fmt.Errorf(MSG_FROM_REQUIRED)
	}

	model, err := s.GetModel(from)
	if err != nil {
		return et.Items{}, err
	}

	data := query.Json("data")
	if len(data) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	where := query.ArrayJson("where")
	if len(where) == 0 {
		return et.Items{}, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	cmd := model.Upsert(data)
	cmd.WhereByJson(where)
	result, err := cmd.Exec()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}
