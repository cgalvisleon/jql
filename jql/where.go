package jql

import (
	"github.com/cgalvisleon/et/et"
)

/**
* Wheres
**/
type Wheres struct {
	owner      *From        `json:"-"`
	conditions []*Condition `json:"-"`
	isDebug    bool         `json:"-"`
}

/**
* newWhere
* @return *Wheres
**/
func newWhere() *Wheres {
	return &Wheres{
		owner:      &From{},
		conditions: make([]*Condition, 0),
	}
}

/**
* setModel
* @param model *Model
* @return *Wheres
**/
func (s *Wheres) setModel(model *Model) *Wheres {
	if model == nil {
		return s
	}
	s.owner = model.from()
	return s
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
			result.Add(condition)
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
	if len(s.conditions) > 0 && condition.Connector == NaC {
		condition.Connector = And
	}
	condition.Field.From = s.owner
	s.conditions = append(s.conditions, condition)
	return s
}

/**
* Where
* @param condition *Condition
* @return *Wheres
**/
func (s *Wheres) Where(condition *Condition) *Wheres {
	condition.Connector = NaC
	return s.add(condition)
}

/**
* And
* @param condition *Condition
* @return *Wheres
**/
func (s *Wheres) And(condition *Condition) *Wheres {
	condition.Connector = And
	return s.add(condition)
}

/**
* Or
* @param condition *Condition
* @return *Wheres
**/
func (s *Wheres) Or(condition *Condition) *Wheres {
	condition.Connector = Or
	return s.add(condition)
}
