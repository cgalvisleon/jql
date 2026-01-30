package jql

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

type Aggregation string

func (s Aggregation) Str() string {
	return string(s)
}

/**
* GetAggregation
* @param tp string
* @return Aggregation
**/
func GetAggregation(tp string) Aggregation {
	aggregation := map[string]Aggregation{
		"count": COUNT,
		"sum":   SUM,
		"avg":   AVG,
		"max":   MAX,
		"min":   MIN,
		"exp":   EXP,
	}

	result, ok := aggregation[tp]
	if !ok {
		return EXP
	}
	return result
}

const (
	COUNT Aggregation = "count"
	SUM   Aggregation = "sum"
	AVG   Aggregation = "avg"
	MAX   Aggregation = "max"
	MIN   Aggregation = "min"
	EXP   Aggregation = "exp"
)

type Status string

const (
	Active    Status = "active"
	Archived  Status = "archived"
	Canceled  Status = "canceled"
	OfSystem  Status = "of_system"
	ForDelete Status = "for_delete"
	Pending   Status = "pending"
	Approved  Status = "approved"
	Rejected  Status = "rejected"
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
func (s *Column) Field() Field {
	return Field{
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
