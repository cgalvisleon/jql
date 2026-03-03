package jql

import (
	"github.com/cgalvisleon/jql/jdb"
)

/**
* NewModel
* @param db *DB, schema, name string, version int
* @return (*Model, error)
**/
func NewModel(db *DB, schema, name string, version int) (*Model, error) {
	result, err := db.NewModel(schema, name, version)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* GetModel
* @param database, name string
* @return *jdb.Model, error
**/
func GetModel(database, name string) (*jdb.Model, error) {
	db, err := jdb.GetDb(database)
	if err != nil {
		return nil, err
	}

	result, err := db.GetModel(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* DeleteModel
* @param database, schema, name string
* @return error
**/
func DeleteModel(database, schema, name string) error {
	db, err := jdb.GetDb(database)
	if err != nil {
		return err
	}

	err = db.DeleteModel(schema, name)
	if err != nil {
		return err
	}

	return nil
}
