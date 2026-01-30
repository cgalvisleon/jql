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
	Database string   `json:"database"`
	Schema   string   `json:"schema"`
	Name     string   `json:"name"`
	As       string   `json:"as"`
	Fields   []*Field `json:"fields"`
}

type Field struct {
	TypeColumn TypeColumn  `json:"type_column"`
	From       *From       `json:"from"`
	Name       interface{} `json:"name"`
	As         string      `json:"as"`
}

func (s *Field) AS() string {
	return fmt.Sprintf("%s.%s", s.From.As, s.As)
}

type Agg struct {
	Agg    string   `json:"agg"`
	Fields []*Field `json:"fields"`
}

type Fld interface {
	string | *Field | *Agg
}

/**
* field
* @param f T
* @return *Field
**/
func field[T Fld](f T) *Field {
	switch v := any(f).(type) {
	case string:
		return &Field{
			Name: v,
			As:   v,
		}
	case *Field:
		return v
	case *Agg:
		return &Field{
			Name: v,
			As:   v.Agg,
		}
	default:
		return nil
	}
}

/**
* agg
* @param agg string, fields ...Fld
* @return *Agg
**/
func agg[T Fld](agg string, fields ...T) *Agg {
	result := &Agg{
		Agg:    agg,
		Fields: make([]*Field, 0),
	}

	for _, f := range fields {
		result.Fields = append(result.Fields, field(f))
	}

	return result
}

/**
* COUNT
* @param fields ...Fld
* @return *Agg
**/
func COUNT[T Fld](fields ...T) *Agg {
	return agg("count", fields...)
}

/**
* SUM
* @param fields ...Fld
* @return *Agg
**/
func SUM[T Fld](fields ...T) *Agg {
	return agg("sum", fields...)
}

/**
* AVG
* @param fields ...Fld
* @return *Agg
**/
func AVG[T Fld](fields ...T) *Agg {
	return agg("avg", fields...)
}

/**
* MAX
* @param fields ...Fld
* @return *Agg
**/
func MAX[T Fld](fields ...T) *Agg {
	return agg("max", fields...)
}

/**
* MIN
* @param fields ...Fld
* @return *Agg
**/
func MIN[T Fld](fields ...T) *Agg {
	return agg("min", fields...)
}

/**
* MIN
* @param fields ...Fld
* @return *Agg
**/
func EXP[T Fld](fields ...T) *Agg {
	return agg("exp", fields...)
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
		Name:       s.Name,
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
