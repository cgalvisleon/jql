package jql

import "fmt"

const (
	SOURCE     string = "source"
	ID         string = "id"
	IDX        string = "idx"
	STATUS     string = "status"
	VERSION    string = "version"
	PROJECT_ID string = "project_id"
	TENANT_ID  string = "tenant_id"
	CREATED_AT string = "created_at"
	UPDATED_AT string = "updated_at"
)

type TypeColumn string

func (s TypeColumn) Str() string {
	return string(s)
}

const (
	COLUMN   TypeColumn = "column"
	ATTRIB   TypeColumn = "atrib"
	DETAIL   TypeColumn = "detail"
	ROLLUP   TypeColumn = "rollup"
	RELATION TypeColumn = "relation"
	CALC     TypeColumn = "calc"
	AGG      TypeColumn = "agg"
)

type TypeData string

func (s TypeData) Str() string {
	return string(s)
}

const (
	ANY      TypeData = "any"
	BYTES    TypeData = "bytes"
	INT      TypeData = "int"
	FLOAT    TypeData = "float"
	KEY      TypeData = "key"
	TEXT     TypeData = "text"
	MEMO     TypeData = "memo"
	JSON     TypeData = "json"
	DATETIME TypeData = "datetime"
	BOOLEAN  TypeData = "boolean"
	GEOMETRY TypeData = "geometry"
)

type From struct {
	Database   string             `json:"database"`
	Schema     string             `json:"schema"`
	Name       string             `json:"name"`
	As         string             `json:"as"`
	Fields     []*Field           `json:"fields"`
	References map[string]*Detail `json:"references"`
	Details    map[string]*Detail `json:"details"`
	Rollups    map[string]*Detail `json:"rollups"`
	Relations  map[string]*Detail `json:"relations"`
}

/**
* findField
* @param name string
* @return *Field
**/
func (s *From) findField(name string) *Field {
	for _, fld := range s.Fields {
		if fld.Field == name {
			return fld
		}
	}
	return nil
}

type Field struct {
	TypeColumn TypeColumn  `json:"type_column"`
	From       *From       `json:"from"`
	Field      interface{} `json:"field"`
	As         string      `json:"as"`
	Page       int         `json:"page"`
	Rows       int         `json:"rows"`
}

/**
* Name
* @return string
**/
func (s *Field) Name() string {
	switch v := s.Field.(type) {
	case string:
		return v
	case *Agg:
		return v.Field
	default:
		return ""
	}
}

/**
* AS
* @return string
**/
func (s *Field) AS() string {
	switch v := s.Field.(type) {
	case string:
		return fmt.Sprintf(`%s.%s:%s`, s.From.As, v, s.As)
	case *Agg:
		return fmt.Sprintf(`%s:%s`, v.AS(), s.As)
	default:
		return fmt.Sprintf(`%v:%s`, v, s.As)
	}
}

type Agg struct {
	Agg   string `json:"agg"`
	Field string `json:"field"`
}

/**
* AS
* @return string
**/
func (s *Agg) AS() string {
	return fmt.Sprintf(`%s(%s):%s`, s.Agg, s.Field, s.Agg)
}

var Aggs = []string{"count", "sum", "avg", "max", "min", "exp"}

/**
* agg
* @param agg string, field string
* @return string
**/
func agg(agg string, field string) *Agg {
	return &Agg{
		Agg:   agg,
		Field: field,
	}
}

/**
* COUNT
* @param field string
* @return *Agg
**/
func COUNT(field string) *Agg {
	return agg("count", field)
}

/**
* SUM
* @param field string
* @return *Agg
**/
func SUM(field string) *Agg {
	return agg("sum", field)
}

/**
* AVG
* @param field string
* @return *Agg
**/
func AVG(field string) *Agg {
	return agg("avg", field)
}

/**
* MAX
* @param field string
* @return *Agg
**/
func MAX(field string) *Agg {
	return agg("max", field)
}

/**
* MIN
* @param field string
* @return *Agg
**/
func MIN(field string) *Agg {
	return agg("min", field)
}

/**
* EXP
* @param field string
* @return *Agg
**/
func EXP(field string) *Agg {
	return agg("exp", field)
}

type Status string

const (
	ACTIVE     Status = "active"
	ARCHIVED   Status = "archived"
	CANCELED   Status = "canceled"
	OF_SYSTEM  Status = "of_system"
	FOR_DELETE Status = "for_delete"
	PENDING    Status = "pending"
	APPROVED   Status = "approved"
	REJECTED   Status = "rejected"
)

type Column struct {
	Name       string      `json:"name"`
	TypeColumn TypeColumn  `json:"type_column"`
	TypeData   TypeData    `json:"type_data"`
	Default    interface{} `json:"default"`
	Definition []byte      `json:"definition"`
	model      *Model      `json:"-"`
}

/**
* Field
* @return Field
**/
func (s *Column) Field() *Field {
	return &Field{
		TypeColumn: s.TypeColumn,
		From:       s.model.from(),
		Field:      s.Name,
		As:         s.Name,
	}
}

/**
* newColumn
* @param model *Model, name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte
* @return *Column
**/
func newColumn(model *Model, name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte) *Column {
	return &Column{
		Name:       name,
		TypeColumn: tpColumn,
		TypeData:   tpData,
		Default:    defaultValue,
		Definition: definition,
		model:      model,
	}
}
