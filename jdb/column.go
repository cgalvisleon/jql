package jdb

const (
	SOURCE     string = "source"
	ID         string = "id"
	IDX        string = "idx"
	STATUS     string = "status"
	VERSION    string = "version"
	TENANT_ID  string = "tenant_id"
	PROJECT_ID string = "project_id"
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

const (
	ACTIVE     string = "active"
	ARCHIVED   string = "archived"
	CANCELED   string = "canceled"
	OF_SYSTEM  string = "of_system"
	FOR_DELETE string = "for_delete"
	PENDING    string = "pending"
	APPROVED   string = "approved"
	REJECTED   string = "rejected"
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
