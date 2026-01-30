package jql

import (
	"github.com/cgalvisleon/et/et"
)

/**
* Wheres
**/
type Wheres struct {
	conditions []*Condition `json:"-"`
	isDebug    bool         `json:"-"`
}

/**
* newWhere
* @return *Wheres
**/
func newWhere() *Wheres {
	return &Wheres{
		conditions: make([]*Condition, 0),
	}
}

/**
* ByJson
* @param jsons []et.Json
* @return *Wheres
**/
func ByJson(jsons []et.Json) *Wheres {
	result := newWhere()
	for _, where := range jsons {
		condition := ToCondition(where)
		if condition != nil {
			result.add(condition)
		}
	}
	return result
}

/**
* IsDebug: Returns the debug mode
* @return *Wheres
**/
func (s *Wheres) IsDebug() *Wheres {
	s.isDebug = true
	return s
}

/**
* ToJson
* @return []et.Json
**/
func (s *Wheres) ToJson() []et.Json {
	result := []et.Json{}
	for _, condition := range s.conditions {
		result = append(result, condition.ToJson())
	}

	return result
}

/**
* add
* @param condition *Condition
* @return *Wheres
**/
func (s *Wheres) add(condition *Condition) *Wheres {
	if len(s.conditions) > 0 && condition.Connector == NAC {
		condition.Connector = AND
	}
	s.conditions = append(s.conditions, condition)
	return s
}
