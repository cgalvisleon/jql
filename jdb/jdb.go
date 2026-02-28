package jdb

import (
	"errors"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

var (
	dbs    map[string]*DB
	models map[string]*Model
)

func init() {
	dbs = make(map[string]*DB)
	models = make(map[string]*Model)
}

/**
* GetDb
* @param name string
* @return *DB, error
**/
func GetDb(name string) (*DB, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "name")
	}

	name = utility.Normalize(name)
	result, ok := dbs[name]
	if ok {
		return result, nil
	}

	exists, err := getCatalog("db", name, result)
	if err != nil {
		return nil, err
	}

	if exists {
		err = result.init()
		if err != nil {
			return nil, err
		}

		dbs[name] = result
		return result, nil
	}

	return nil, errors.New(MSG_DB_NOT_FOUND)
}

/**
* NewDb
* @param name string, params et.Json
* @return *DB, error
**/
func NewDb(name string, params et.Json) (*DB, error) {
	result, err := GetDb(name)
	if err != nil {
		return nil, err
	}

	if result != nil {
		return result, nil
	}

	driver := params.Str("driver")
	drv, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
	}

	result = &DB{
		Name:    name,
		Schemas: make(map[string]*Schema),
		Params:  params,
		driver:  drv(),
	}
	err = result.init()
	if err != nil {
		return nil, err
	}

	dbs[name] = result
	return result, nil
}

/**
* DeleteDb
* @param name string
* @return error
**/
func DeleteDb(name string) error {
	_, ok := dbs[name]
	if ok {
		delete(dbs, name)
	}

	err := deleteCatalog("db", name)
	if err != nil {
		return err
	}

	return nil
}
