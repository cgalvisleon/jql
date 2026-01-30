package jql

import (
	"encoding/json"
	"regexp"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/timezone"
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
	Unique        []string               `json:"unique"`
	Required      []string               `json:"required"`
	Hidden        []string               `json:"hidden"`
	Details       map[string]*Detail     `json:"details"`
	Rollups       map[string]*Detail     `json:"rollups"`
	Relations     map[string]*Detail     `json:"relations"`
	BeforeInserts []*Trigger             `json:"before_inserts"`
	BeforeUpdates []*Trigger             `json:"before_updates"`
	BeforeDeletes []*Trigger             `json:"before_deletes"`
	AfterInserts  []*Trigger             `json:"after_inserts"`
	AfterUpdates  []*Trigger             `json:"after_updates"`
	AfterDeletes  []*Trigger             `json:"after_deletes"`
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
* save
* @return error
**/
func (s *Model) save() error {
	if models == nil {
		return nil
	}

	if s.IsCore {
		return nil
	}

	serialize, err := s.serialize()
	if err != nil {
		return err
	}

	now := timezone.Now()
	_, err = models.
		Upsert(et.Json{
			"name":       s.Name,
			"version":    s.Version,
			"definition": serialize,
		}).
		BeforeInsertOrUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("created_at", now)
			new.Set("updated_at", now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("updated_at", now)
			return nil
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if s.isInit {
		return nil
	}

	s.isInit = true
	return s.save()
}

/**
* Stricted
* @return error
**/
func (s *Model) Stricted() {
	s.IsStrict = true
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
		As:       s.Name,
		Fields:   make([]*Field, 0),
	}
	for _, column := range s.Columns {
		result.Fields = append(result.Fields, column.Field())
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
* FindField
* @param name string
* @return *Field
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
* findField
* @param name string
* @return *Field
**/
func (s *Model) findField(name string) *Field {
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
* Select
* @param fields ...interface{}
* @return *Ql
**/
func (s *Model) Select(fields ...interface{}) *Ql {
	result := newQuery(s, "A")
	result.Select(fields...)
	return result
}

/**
* Counted
* @return *Ql
**/
func (s *Model) Counted() *Ql {
	result := newQuery(s, "A")
	result.Type = COUNTED
	return result
}

/**
* ItExists
* @return *Ql
**/
func (s *Model) ItExists() *Ql {
	result := newQuery(s, "A")
	result.Type = EXISTS
	return result
}

/**
* Current
* @return *Ql
**/
func (s *Model) Current(data et.Json) *Ql {
	result := newQuery(s, "A")
	for _, col := range s.Columns {
		if col.TypeColumn == COLUMN {
			field := col.Field()
			result.Selects = append(result.Selects, field)
		}
	}
	for _, key := range s.PrimaryKeys {
		if _, ok := data[key]; !ok {
			continue
		}
		result.Wheres.add(Eq(key, data[key]))
	}
	return result
}
