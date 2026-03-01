package jdb

import (
	"encoding/json"
	"errors"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var catalog *Model

/**
* defineCatalog
* @param db *DB
* @return error
**/
func defineCatalog(db *DB) error {
	if catalog != nil {
		return nil
	}

	var err error
	catalog, err = db.NewModel("core", "catalog", 1)
	if err != nil {
		return err
	}
	catalog.DefineCreatedAtField()
	catalog.DefineUpdatedAtField()
	catalog.DefineColumn("type", TEXT, "")
	catalog.DefineColumn("name", KEY, "")
	catalog.DefineColumn("version", INT, 0)
	catalog.DefineColumn("definition", BYTES, []byte{})
	catalog.DefinePrimaryKeys("type", "name")
	catalog.IsCore = true
	if err = catalog.Init(); err != nil {
		return err
	}

	return nil
}

/**
* setCatalog
* @param tp, name string, version int, obj any
* @return error
**/
func setCatalog(tp, name string, version int, obj any) error {
	if catalog == nil {
		return nil
	}

	bt, ok := obj.([]byte)
	if !ok {
		var err error
		bt, err = json.Marshal(obj)
		if err != nil {
			return err
		}
	}

	now := timezone.Now()
	_, err := catalog.
		Upsert(et.Json{
			"type":       tp,
			"name":       name,
			"version":    version,
			"definition": bt,
		}).
		BeforeInsertOrUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("created_at", now)
			new.Set("updated_at", now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, old, new et.Json) error {
			oldVersion := old.Int("version")
			if oldVersion == version {
				return ErrNotUpdated
			}

			new.Set("updated_at", now)
			return nil
		}).
		Exec()
	if err != nil && !errors.Is(err, ErrNotUpdated) {
		return err
	}

	return nil
}

/**
* getCatalog
* @param tp, name string, dest any
* @return bool, error
**/
func getCatalog(tp, name string, dest any) (bool, error) {
	item, err := catalog.
		Where(Eq("type", tp)).
		And(Eq("name", name)).
		Select().
		One()
	if err != nil {
		return false, err
	}

	if !item.Ok {
		return false, nil
	}

	bt, err := item.Byte("definition")
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(bt, dest)
	if err != nil {
		return false, err
	}

	return true, nil
}

/**
* deleteCatalog
* @param tp, name string
* @return error
**/
func deleteCatalog(tp, name string) error {
	_, err := catalog.
		Delete().
		Where(Eq("type", tp)).
		And(Eq("name", name)).
		One()
	if err != nil {
		return err
	}

	return nil
}
