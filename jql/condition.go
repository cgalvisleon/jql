package jql

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Operator string

const (
	OpEq         Operator = "eq"
	OpNeg        Operator = "neg"
	OpLess       Operator = "less"
	OpLessEq     Operator = "less_eq"
	OpMore       Operator = "more"
	OpMoreEq     Operator = "more_eq"
	OpLike       Operator = "like"
	OpIn         Operator = "in"
	OpNotIn      Operator = "not_in"
	OpIs         Operator = "is"
	OpIsNot      Operator = "is_not"
	OpNull       Operator = "null"
	OpNotNull    Operator = "not_null"
	OpBetween    Operator = "between"
	OpNotBetween Operator = "not_between"
)

func (s Operator) Str() string {
	return string(s)
}

func ToOperator(s string) Operator {
	values := map[string]Operator{
		"eq":          OpEq,
		"neg":         OpNeg,
		"less":        OpLess,
		"less_eq":     OpLessEq,
		"more":        OpMore,
		"more_eq":     OpMoreEq,
		"like":        OpLike,
		"in":          OpIn,
		"not_in":      OpNotIn,
		"is":          OpIs,
		"is_not":      OpIsNot,
		"null":        OpNull,
		"not_null":    OpNotNull,
		"between":     OpBetween,
		"not_between": OpNotBetween,
	}

	result, ok := values[s]
	if !ok {
		return OpEq
	}

	return result
}

type Connector string

const (
	NAC Connector = ""
	AND Connector = "and"
	OR  Connector = "or"
)

func (s Connector) Str() string {
	return string(s)
}

type BetweenValue struct {
	Min any
	Max any
}

type Condition struct {
	Field     *Field    `json:"field"`
	Operator  Operator  `json:"operator"`
	Value     any       `json:"value"`
	Connector Connector `json:"connector"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Condition) ToJson() et.Json {
	if s.Connector == NAC {
		return et.Json{
			s.Field.AS(): et.Json{
				s.Operator.Str(): s.Value,
			},
		}
	}

	return et.Json{
		s.Connector.Str(): et.Json{
			s.Field.AS(): et.Json{
				s.Operator.Str(): s.Value,
			},
		},
	}
}

/**
* ToCondition
* @param json et.Json
* @return *Condition
**/
func ToCondition(json et.Json) *Condition {
	getWhere := func(json et.Json) *Condition {
		for fld := range json {
			cond := json.Json(fld)
			for cnd := range cond {
				val := cond[cnd]
				return condition(fld, val, ToOperator(cnd))
			}
		}
		return nil
	}

	and := func(jsons et.Json) *Condition {
		result := getWhere(jsons)
		if result != nil {
			result.Connector = AND
		}

		return result
	}

	or := func(jsons et.Json) *Condition {
		result := getWhere(jsons)
		if result != nil {
			result.Connector = OR
		}

		return result
	}

	for k := range json {
		if strs.Lowcase(k) == "and" {
			def := json.Json(k)
			return and(def)
		} else if strs.Lowcase(k) == "or" {
			def := json.Json(k)
			return or(def)
		} else {
			return getWhere(json)
		}
	}

	return nil
}

/**
* condition
* @param field interface{}, value interface{}, op Operator
* @return *Condition
**/
func condition(field interface{}, value interface{}, op Operator) *Condition {
	return &Condition{
		Field:     &Field{Field: field},
		Operator:  op,
		Value:     value,
		Connector: NAC,
	}
}

/**
* Eq
* @param field interface{}, value interface{}
* @return Condition
**/
func Eq(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpEq)
}

/**
* Neg
* @param field interface{}, value interface{}
* @return Condition
**/
func Neg(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpNeg)
}

/**
* Less
* @param field interface{}, value interface{}
* @return Condition
**/
func Less(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpLess)
}

/**
* LessEq
* @param field interface{}, value interface{}
* @return Condition
**/
func LessEq(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpLessEq)
}

/**
* More
* @param field interface{}, value interface{}
* @return Condition
**/
func More(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpMore)
}

/**
* MoreEq
* @param field interface{}, value interface{}
* @return Condition
**/
func MoreEq(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpMoreEq)
}

/**
* Like
* @param field interface{}, value interface{}
* @return Condition
**/
func Like(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpLike)
}

/**
* In
* @param field interface{}, value []interface{}
* @return Condition
**/
func In(field interface{}, value []interface{}) *Condition {
	return condition(field, value, OpIn)
}

/**
* NotIn
* @param field interface{}, value []interface{}
* @return Condition
**/
func NotIn(field interface{}, value []interface{}) *Condition {
	return condition(field, value, OpNotIn)
}

/**
* Is
* @param field interface{}, value interface{}
* @return Condition
**/
func Is(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpIs)
}

/**
* IsNot
* @param field interface{}, value interface{}
* @return Condition
**/
func IsNot(field interface{}, value interface{}) *Condition {
	return condition(field, value, OpIsNot)
}

/**
* Null
* @param field interface{}
* @return Condition
**/
func Null(field interface{}) *Condition {
	return condition(field, nil, OpNull)
}

/**
* NotNull
* @param field interface{}
* @return Condition
**/
func NotNull(field interface{}) *Condition {
	return condition(field, nil, OpNotNull)
}

/**
* Between
* @param field interface{}, min any, max any
* @return Condition
**/
func Between(field interface{}, min, max any) *Condition {
	return condition(field, BetweenValue{Min: min, Max: max}, OpBetween)
}

/**
* NotBetween
* @param field interface{}, min any, max any
* @return Condition
**/
func NotBetween(field interface{}, min, max any) *Condition {
	return condition(field, BetweenValue{Min: min, Max: max}, OpNotBetween)
}
