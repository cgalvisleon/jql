package jql

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineColumn
* @param name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte
* @return *Column
**/
func (s *Model) defineColumn(name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte) (*Column, error) {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	if !utility.ValidStr(tpColumn.Str(), 0, []string{}) {
		return nil, fmt.Errorf(MSG_TYPE_COLUMN_REQUIRED)
	}

	if !utility.ValidStr(tpData.Str(), 0, []string{}) {
		return nil, fmt.Errorf(MSG_TYPE_DATA_REQUIRED)
	}

	idx := s.idxColumn(name)
	if idx != -1 {
		return s.Columns[idx], nil
	}

	result := newColumn(s, name, tpColumn, tpData, defaultValue, definition)
	s.Columns = append(s.Columns, result)
	return result, nil
}

/**
* DefineIndex
* @param names ...string
**/
func (s *Model) DefineIndex(names ...string) {
	for _, name := range names {
		idx := s.idxColumn(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.Indexes, name)
		if idx != -1 {
			continue
		}

		s.Indexes = append(s.Indexes, name)
	}
}

/**
* DefinePrimaryKeys
* @param names ...string
**/
func (s *Model) DefinePrimaryKeys(names ...string) {
	for _, name := range names {
		idx := s.idxColumn(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.PrimaryKeys, name)
		if idx != -1 {
			continue
		}

		s.PrimaryKeys = append(s.PrimaryKeys, name)
		s.DefineRequired(name)
	}
}

/**
* DefineUnique
* @param names ...string
**/
func (s *Model) DefineUnique(names ...string) {
	for _, name := range names {
		idx := s.idxColumn(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.Unique, name)
		if idx != -1 {
			continue
		}

		s.Unique = append(s.Unique, name)
	}
}

/**
* DefineRequired
* @param names ...string
**/
func (s *Model) DefineRequired(names ...string) {
	for _, name := range names {
		idx := s.idxColumn(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.Required, name)
		if idx != -1 {
			continue
		}

		s.Required = append(s.Required, name)
	}
}

/**
* DefineHidden
* @param names ...string
**/
func (s *Model) DefineHidden(names ...string) {
	for _, name := range names {
		idx := s.idxColumn(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.Hidden, name)
		if idx != -1 {
			continue
		}

		s.Hidden = append(s.Hidden, name)
	}
}

/**
* DefineColumn
* @param name string, tpData TypeData, defaultValue interface{}
* @return *Column
**/
func (s *Model) DefineColumn(name string, tpData TypeData, defaultValue interface{}) (*Column, error) {
	return s.defineColumn(name, COLUMN, tpData, defaultValue, []byte{})
}

/**
* DefineSourceField
* @param name string
* @return error
**/
func (s *Model) DefineSourceField(name string) (*Column, error) {
	result, err := s.DefineColumn(name, JSON, et.Json{})
	if err != nil {
		return nil, err
	}

	s.SourceField = name
	s.DefineIndex(name)
	return result, nil
}

/**
* DefineIdxField
* @param name string
* @return error
**/
func (s *Model) DefineIdxField(name string) (*Column, error) {
	result, err := s.DefineColumn(name, KEY, "")
	if err != nil {
		return nil, err
	}

	s.IdxField = name
	s.DefineIndex(name)
	return result, nil
}

/**
* DefineAttribute
* @param name string, tpData TypeData, defaultValue interface{}
* @return *Column, error
**/
func (s *Model) DefineAttribute(name string, tpData TypeData, defaultValue interface{}) (*Column, error) {
	if s.SourceField == "" {
		_, err := s.DefineSourceField(SOURCE)
		if err != nil {
			return nil, err
		}
	}
	return s.defineColumn(name, ATTRIB, tpData, defaultValue, []byte{})
}

/**
* DefineDetail
* @param name string, to *Model, keys map[string]string
* @return *Column
**/
func (s *Model) DefineDetail(name string, keys map[string]string, version int) (*Model, error) {
	_, err := s.defineColumn(name, DETAIL, JSON, []et.Json{}, []byte{})
	if err != nil {
		return nil, err
	}

	to, err := s.db.newModel(s.Schema, fmt.Sprintf("%s_%s", s.Name, name), version)
	if err != nil {
		return nil, err
	}

	for fk, pk := range keys {
		_, err = s.DefineColumn(pk, KEY, "")
		if err != nil {
			return nil, err
		}

		_, err = to.DefineColumn(fk, KEY, "")
		if err != nil {
			return nil, err
		}
	}

	detail := newDetail(to, keys, []string{}, true, true)
	s.Details[name] = detail
	return to, nil
}

/**
* DefineRollup
* @param name string, from string, keys map[string]string, selects []string
* @return *Model
**/
func (s *Model) DefineRollup(name string, from *Model, keys map[string]string, selects []string) error {
	_, err := s.defineColumn(name, ROLLUP, JSON, "", []byte{})
	if err != nil {
		return err
	}

	detail := newDetail(from, keys, selects, false, false)
	s.Rollups[name] = detail
	return nil
}

/**
* DefineRelation
* @param name string, from string, keys map[string]string
* @return *Model
**/
func (s *Model) DefineRelation(from *Model, keys map[string]string) error {
	detail := newDetail(from, keys, []string{}, false, false)
	s.Relations[from.Name] = detail
	return nil
}

/**
* DefineCalc
* @param name string, fn DataContext
* @return error
**/
func (s *Model) DefineCalc(name string, fn DataContext) error {
	_, err := s.defineColumn(name, TpCalc, TpAny, false, nil, []byte{})
	if err != nil {
		return err
	}

	s.calcs[name] = fn
	return nil
}

/**
* DefineCreatedAtField
* @return *Model
**/
func (s *Model) DefineCreatedAtField() *Model {
	s.DefineColumn(CREATED_AT, TpDateTime)
	return s
}

/**
* DefineUpdatedAtField
* @return *Model
**/
func (s *Model) DefineUpdatedAtField() *Model {
	s.DefineColumn(UPDATED_AT, TpDateTime)
	return s
}

/**
* DefineStatusFieldDefault
* @return *Model
**/
func (s *Model) DefineStatusFieldDefault() *Model {
	s.DefineColumn(STATUS, TpKey)
	return s
}

/**
* DefinePrimaryKeyField
* @return *Model
**/
func (s *Model) DefinePrimaryKeyField() *Model {
	s.DefinePrimaryKeys(KEY)
	return s
}

/**
* DefineSourceFieldDefault
* @return *Model
**/
func (s *Model) DefineSourceFieldDefault() *Model {
	s.DefineSourceField(SOURCE)
	return s
}

func (s *Model) DefineIndexFieldDefault() *Model {
	s.DefineIndexField(IDX)
	return s
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusFieldDefault()
	s.DefinePrimaryKeyField()
	s.DefineSourceFieldDefault()
	s.DefineIndexFieldDefault()
	return s
}

/**
* DefineTenantModel
* @return *Model
**/
func (s *Model) DefineProjectModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusFieldDefault()
	s.DefinePrimaryKeyField()
	s.DefineColumn(PROJECT_ID, TpKey)
	s.DefineSourceFieldDefault()
	s.DefineIndexFieldDefault()
	s.DefineIndex(PROJECT_ID)
	return s
}
