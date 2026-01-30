package jql

import (
	"errors"

	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/josefina/pkg/msg"
)

type Schema struct {
	Database string            `json:"-"`
	Name     string            `json:"name"`
	Models   map[string]*Model `json:"models"`
}

func (s *Schema) newModel(name string, isCore bool, version int) (*Model, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return nil, errors.New(msg.MSG_MODEL_NOT_FOUND)
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
		Unique:        make([]string, 0),
		Required:      make([]string, 0),
		Hidden:        make([]string, 0),
		References:    make(map[string]*Detail, 0),
		Details:       make(map[string]*Detail, 0),
		Rollups:       make(map[string]*Detail, 0),
		Relations:     make(map[string]*Detail, 0),
		Calcs:         make(map[string][]byte, 0),
		BeforeInserts: make([]*Trigger, 0),
		BeforeUpdates: make([]*Trigger, 0),
		BeforeDeletes: make([]*Trigger, 0),
		AfterInserts:  make([]*Trigger, 0),
		AfterUpdates:  make([]*Trigger, 0),
		AfterDeletes:  make([]*Trigger, 0),
		IsCore:        isCore,
		Version:       version,
		beforeInserts: make([]TriggerFunction, 0),
		beforeUpdates: make([]TriggerFunction, 0),
		beforeDeletes: make([]TriggerFunction, 0),
		afterInserts:  make([]TriggerFunction, 0),
		afterUpdates:  make([]TriggerFunction, 0),
		afterDeletes:  make([]TriggerFunction, 0),
		calcs:         make(map[string]DataContext),
	}
	s.Models[name] = result

	return result, nil
}
