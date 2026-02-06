package jql

import (
	"fmt"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Schema struct {
	Database string            `json:"-"`
	Name     string            `json:"name"`
	Models   map[string]*From `json:"models"`
	db       *DB               `json:"-"`
}

/**
* newModel
* @param name string, version int
* @return (*Model, error)
**/
func (s *Schema) newModel(name string, version int) (*Model, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "name")
	}

	result, ok := s.Models[name]
	if ok {
		return result, nil
	}

	name = utility.Normalize(name)
	result = &Model{
		Database:      s.Database,
		Schema:        s.Name,
		Name:          name,
		Columns:       make([]*Column, 0),
		Indexes:       make([]string, 0),
		PrimaryKeys:   make([]string, 0),
		ForeignKeys:   make([]*Detail, 0),
		Unique:        make([]string, 0),
		Required:      make([]string, 0),
		Hidden:        make([]string, 0),
		Details:       make(map[string]*Detail, 0),
		Rollups:       make(map[string]*Detail, 0),
		Relations:     make(map[string]*Detail, 0),
		BeforeInserts: make([]*Trigger, 0),
		BeforeUpdates: make([]*Trigger, 0),
		BeforeDeletes: make([]*Trigger, 0),
		AfterInserts:  make([]*Trigger, 0),
		AfterUpdates:  make([]*Trigger, 0),
		AfterDeletes:  make([]*Trigger, 0),
		Version:       version,
		beforeInserts: make([]TriggerFunction, 0),
		beforeUpdates: make([]TriggerFunction, 0),
		beforeDeletes: make([]TriggerFunction, 0),
		afterInserts:  make([]TriggerFunction, 0),
		afterUpdates:  make([]TriggerFunction, 0),
		afterDeletes:  make([]TriggerFunction, 0),
		calcs:         make(map[string]DataContext),
	}
	result.defineIdxField()
	s.Models[name] = result

	return result, nil
}

/**
* getModel
* @param name string
* @return (*Model, error)
**/
func (s *Schema) getModel(name string) (*Model, error) {
	result, ok := s.Models[name]
	if ok {
		return result, nil
	}

	if catalog == nil {
		return nil, ErrModelNotFound
	}

	key := name
	key = strs.Append(s.Name, key, ".")
	key = strs.Append(s.Database, key, ".")
	exists, err := getCatalog("model", key, &result)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrModelNotFound
	}

	result.beforeInserts = make([]TriggerFunction, 0)
	result.beforeUpdates = make([]TriggerFunction, 0)
	result.beforeDeletes = make([]TriggerFunction, 0)
	result.afterInserts = make([]TriggerFunction, 0)
	result.afterUpdates = make([]TriggerFunction, 0)
	result.afterDeletes = make([]TriggerFunction, 0)
	result.calcs = make(map[string]DataContext)
	result.db = s.db
	err = result.Init()
	if err != nil {
		return nil, err
	}

	s.Models[name] = result
	return result, nil
}
