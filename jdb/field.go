package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/strs"
)

type From struct {
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	Table    string `json:"table"`
	As       string `json:"as"`
	model    *Model `json:"-"`
}

func (s *From) Key() string {
	result := s.model.Key()
	return result
}

type Agg struct {
	Agg   string `json:"agg"`
	Field string `json:"field"`
}

/**
* Name
* @return string
**/
func (s *Agg) Name() string {
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
		result := v
		result = strs.Append(result, s.As, ":")
		return result
	case *Agg:
		return v.Name()
	default:
		return ""
	}
}
