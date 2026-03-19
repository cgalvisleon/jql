package jdb

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
)

type Trigger struct {
	Name       string `json:"name"`
	Definition []byte `json:"definition"`
}

type TriggerFunction func(tx *Tx, old, new et.Json) error

type DataContext func(tx *Tx, data et.Json)

type Model struct {
	Database      string                 `json:"database"`
	Schema        string                 `json:"schema"`
	Name          string                 `json:"name"`
	Table         string                 `json:"table"`
	Columns       []*Column              `json:"columns"`
	SourceField   string                 `json:"source_field"`
	IdxField      string                 `json:"idx_field"`
	Indexes       []string               `json:"indexes"`
	PrimaryKeys   []string               `json:"primary_keys"`
	ForeignKeys   []*Detail              `json:"foreign_keys"`
	Unique        []string               `json:"unique"`
	Required      []string               `json:"required"`
	Hidden        []string               `json:"hidden"`
	Details       map[string]*Detail     `json:"details"`
	Rollups       map[string]*Detail     `json:"rollups"`
	Relations     map[string]*Detail     `json:"relations"`
	IsStrict      bool                   `json:"is_strict"`
	Version       int                    `json:"version"`
	IsCore        bool                   `json:"is_core"`
	IsDebug       bool                   `json:"-"`
	isInit        bool                   `json:"-"`
	beforeInserts []TriggerFunction      `json:"-"`
	beforeUpdates []TriggerFunction      `json:"-"`
	beforeDeletes []TriggerFunction      `json:"-"`
	afterInserts  []TriggerFunction      `json:"-"`
	afterUpdates  []TriggerFunction      `json:"-"`
	afterDeletes  []TriggerFunction      `json:"-"`
	calcs         map[string]DataContext `json:"-"`
	db            *DB                    `json:"-"`
}

/**
* serialize
* @return []byte, error
**/
func (s *Model) serialize() ([]byte, error) {
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
func (s *Model) ToJson() et.Json {
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
* Key
* @return string
**/
func (s *Model) Key() string {
	result := s.Name
	result = strs.Append(s.Schema, result, ".")
	result = strs.Append(s.Database, result, ".")
	return result
}

/**
* Save
* @return error
**/
func (s *Model) Save() error {
	if s.IsCore {
		return nil
	}

	serialize, err := s.serialize()
	if err != nil {
		return err
	}

	key := s.Key()
	return setCatalog("model", key, s.Version, serialize)
}

/**
* Debug
**/
func (s *Model) Debug() {
	s.IsDebug = true
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if s.isInit {
		return nil
	}

	err := s.db.loadModel(s)
	if err != nil {
		return err
	}

	s.isInit = true
	if s.IsCore {
		return nil
	}

	oldVersion, err := versionCatalog("model", s.Key())
	if err != nil {
		return err
	}

	if oldVersion < s.Version {
		return s.Save()
	}

	return nil
}

/**
* Stricted
* @return error
**/
func (s *Model) Stricted() {
	s.IsStrict = true
}

/**
* Db
* @return *sql.DB
**/
func (s *Model) Db() *sql.DB {
	return s.db.db
}

/**
* from
* @return From
**/
func (s *Model) from() *From {
	result := &From{
		Database: s.Database,
		Schema:   s.Schema,
		Name:     s.Name,
		Table:    s.Table,
		As:       s.Name,
		model:    s,
	}
	return result
}

/**
* idxColumn
* @param name string
* @return int
**/
func (s *Model) idxColumn(name string) int {
	return slices.IndexFunc(s.Columns, func(column *Column) bool { return column.Name == name })
}

/**
* FindColumn
* @param name string
* @return *Column
**/
func (s *Model) FindColumn(name string) *Column {
	idx := s.idxColumn(name)
	if idx != -1 {
		return s.Columns[idx]
	}

	if s.IsStrict {
		return nil
	}

	if s.SourceField == "" {
		return nil
	}

	return newColumn(s, name, ATTRIB, ANY, "", []byte{})
}

/**
* FindField
* @param name string
* @return *Field
**/
func (s *Model) FindField(name string) *Field {
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9>]+):([A-Za-z0-9]+)$`) // name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9>]+)$`)                // name

	if pattern1.MatchString(name) {
		matches := pattern1.FindStringSubmatch(name)
		if len(matches) == 3 {
			name = matches[1]
			as := matches[2]
			column := s.FindColumn(name)
			if column != nil {
				result := column.Field()
				result.As = as
				return result
			}
		}
	} else if pattern2.MatchString(name) {
		column := s.FindColumn(name)
		if column != nil {
			result := column.Field()
			return result
		}
	}

	return nil
}

/**
* GetId
* @param id string
* @return string
**/
func (s *Model) GetId(id string) string {
	return reg.TagULID(s.Name, id)
}

/**
* BeforeInsert
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) BeforeInsert(fn TriggerFunction) *Model {
	s.beforeInserts = append(s.beforeInserts, fn)
	return s
}

/**
* BeforeUpdate
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) BeforeUpdate(fn TriggerFunction) *Model {
	s.beforeUpdates = append(s.beforeUpdates, fn)
	return s
}

/**
* BeforeDelete
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) BeforeDelete(fn TriggerFunction) *Model {
	s.beforeDeletes = append(s.beforeDeletes, fn)
	return s
}

/**
* BeforeInsertOrUpdate
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) BeforeInsertOrUpdate(fn TriggerFunction) *Model {
	s.beforeInserts = append(s.beforeInserts, fn)
	s.beforeUpdates = append(s.beforeUpdates, fn)
	return s
}

/**
* AfterInsert
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) AfterInsert(fn TriggerFunction) *Model {
	s.afterInserts = append(s.afterInserts, fn)
	return s
}

/**
* AfterUpdate
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) AfterUpdate(fn TriggerFunction) *Model {
	s.afterUpdates = append(s.afterUpdates, fn)
	return s
}

/**
* AfterDelete
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) AfterDelete(fn TriggerFunction) *Model {
	s.afterDeletes = append(s.afterDeletes, fn)
	return s
}

/**
* AfterInsertOrUpdate
* @param fn TriggerFunction
* @return *Model
**/
func (s *Model) AfterInsertOrUpdate(fn TriggerFunction) *Model {
	s.afterInserts = append(s.afterInserts, fn)
	s.afterUpdates = append(s.afterUpdates, fn)
	return s
}

/**
* Insert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Insert(data et.Json) *Cmd {
	result := newCommand(s, INSERT)
	result.Data = append(result.Data, data)
	return result
}

/**
* Update
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Update(data et.Json) *Cmd {
	result := newCommand(s, UPDATE)
	result.Data = append(result.Data, data)
	return result
}

/**
* Delete
* @return *Cmd
**/
func (s *Model) Delete() *Cmd {
	result := newCommand(s, DELETE)
	return result
}

/**
* Upsert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Upsert(data et.Json) *Cmd {
	result := newCommand(s, UPSERT)
	result.Data = append(result.Data, data)
	return result
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func (s *Model) Query(query et.Json) (et.Items, error) {
	result := NewQuery(s, "A")
	return result.Query(query)
}
