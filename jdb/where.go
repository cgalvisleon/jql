package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* Wheres
**/
type Wheres struct {
	Conditions []*Condition `json:"Conditions"`
	isDebug    bool         `json:"-"`
}

/**
* newWhere
* @return *Wheres
**/
func newWhere() *Wheres {
	return &Wheres{
		Conditions: make([]*Condition, 0),
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
	for _, condition := range s.Conditions {
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
	if len(s.Conditions) > 0 && condition.Connector == NAC {
		condition.Connector = AND
	}
	s.Conditions = append(s.Conditions, condition)
	return s
}

/**
* ByPk
* @param model *From, data et.Json
* @return *Wheres
**/
func (s *Wheres) ByPk(model interface{}, data et.Json) *Wheres {
	switch v := model.(type) {
	case *From:
		for _, key := range v.model.PrimaryKeys {
			if _, ok := data[key]; !ok {
				continue
			}
			col := v.model.FindColumn(key)
			if col == nil {
				continue
			}
			s.add(Eq(col, data[key]))
		}
	case *Model:
		for _, key := range v.PrimaryKeys {
			if _, ok := data[key]; !ok {
				continue
			}
			col := v.FindColumn(key)
			if col == nil {
				continue
			}
			s.add(Eq(col, data[key]))
		}
	}

	return s
}
