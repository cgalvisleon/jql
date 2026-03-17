package jdb

import (
	"errors"
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineColumn
* @param name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte
* @return *Column
**/
func (s *Model) defineColumn(name string, tpColumn TypeColumn, tpData TypeData, defaultValue interface{}, definition []byte) (*Column, error) {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, errors.New(MSG_NAME_REQUIRED)
	}

	if !utility.ValidStr(tpColumn.Str(), 0, []string{}) {
		return nil, errors.New(MSG_TYPE_COLUMN_REQUIRED)
	}

	if !utility.ValidStr(tpData.Str(), 0, []string{}) {
		return nil, errors.New(MSG_TYPE_DATA_REQUIRED)
	}

	idx := s.idxColumn(name)
	if idx != -1 {
		return s.Columns[idx], nil
	}

	result := newColumn(s, name, tpColumn, tpData, defaultValue, definition)
	s.Columns = insertBeforeLast(s.Columns, result)
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
* DefineTable
* @param name string
* @return error
**/
func (s *Model) DefineForeignKey(to *Model, keys map[string]string, onDeleteCascade, onUpdateCascade bool) error {
	detail := newDetail(to, keys, []interface{}{}, onDeleteCascade, onUpdateCascade)
	for fk, pk := range keys {
		fld := s.findField(pk)
		if fld == nil {
			return fmt.Errorf(MSG_FIELD_NOT_FOUND, pk)
		}

		fld = to.findField(fk)
		if fld == nil {
			return fmt.Errorf(MSG_FIELD_NOT_FOUND, fk)
		}
	}
	s.ForeignKeys = append(s.ForeignKeys, detail)
	return nil
}

/**
* DefineSourceField
* @return error
**/
func (s *Model) DefineSourceField() (*Column, error) {
	result, err := s.DefineColumn(SOURCE, JSON, et.Json{})
	if err != nil {
		return nil, err
	}

	s.SourceField = SOURCE
	s.DefineIndex(SOURCE)
	return result, nil
}

/**
* defineIdxField
* @return error
**/
func (s *Model) defineIdxField() (*Column, error) {
	result, err := s.DefineColumn(IDX, KEY, "")
	if err != nil {
		return nil, err
	}

	s.IdxField = IDX
	s.DefineIndex(IDX)
	s.BeforeInsert(func(tx *Tx, old, new et.Json) error {
		new[IDX] = reg.ULID()
		return nil
	})
	return result, nil
}

/**
* DefineAttribute
* @param name string, tpData TypeData, defaultValue interface{}
* @return *Column, error
**/
func (s *Model) DefineAttribute(name string, tpData TypeData, defaultValue interface{}) (*Column, error) {
	if s.SourceField == "" {
		_, err := s.DefineSourceField()
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

	to, err := s.db.NewModel(s.Schema, fmt.Sprintf("%s_%s", s.Name, name), version)
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

	detail := newDetail(to, keys, []interface{}{}, true, true)
	s.Details[name] = detail
	return to, nil
}

/**
* DefineRollup
* @param name string, from string, keys map[string]string, selects []interface{}
* @return *Model
**/
func (s *Model) DefineRollup(name string, from *Model, keys map[string]string, selects []interface{}) error {
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
	detail := newDetail(from, keys, []interface{}{}, false, false)
	s.Relations[from.Name] = detail
	return nil
}

/**
* DefineCalc
* @param name string, fn DataContext
* @return error
**/
func (s *Model) DefineCalc(name string, fn DataContext) error {
	_, err := s.defineColumn(name, CALC, ANY, []byte{}, []byte{})
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
	s.DefineColumn(CREATED_AT, DATETIME, "")
	return s
}

/**
* DefineUpdatedAtField
* @return *Model
**/
func (s *Model) DefineUpdatedAtField() *Model {
	s.DefineColumn(UPDATED_AT, DATETIME, "")
	return s
}

/**
* DefineStatusField
* @return *Model
**/
func (s *Model) DefineStatusField() *Model {
	s.DefineColumn(STATUS, KEY, "")
	return s
}

/**
* DefinePrimaryKeyField
* @return *Model
**/
func (s *Model) DefinePrimaryKeyField() *Model {
	s.DefineColumn(ID, KEY, "")
	s.DefinePrimaryKeys(ID)
	return s
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusField()
	s.DefinePrimaryKeyField()
	s.DefineSourceField()
	s.BeforeInsert(func(tx *Tx, old, new et.Json) error {
		new.Set(CREATED_AT, timezone.Now())
		new.Set(UPDATED_AT, timezone.Now())
		return nil
	})
	s.BeforeUpdate(func(tx *Tx, old, new et.Json) error {
		new.Set(UPDATED_AT, timezone.Now())
		return nil
	})
	return s
}

/**
* DefineTenantModel
* @return *Model
**/
func (s *Model) DefineProjectModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusField()
	s.DefinePrimaryKeyField()
	s.DefineColumn(PROJECT_ID, KEY, "")
	s.defineIdxField()
	s.DefineIndex(PROJECT_ID)
	s.BeforeInsert(func(tx *Tx, old, new et.Json) error {
		new.Set(CREATED_AT, timezone.Now())
		new.Set(UPDATED_AT, timezone.Now())
		id := reg.GenULID(s.Name)
		new.Set(ID, id)
		return nil
	})
	s.BeforeUpdate(func(tx *Tx, old, new et.Json) error {
		new.Set(UPDATED_AT, timezone.Now())
		return nil
	})
	return s
}
