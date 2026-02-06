package jql

import (
	"encoding/json"

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
	catalog.defineCreatedAtField()
	catalog.defineUpdatedAtField()
	catalog.DefineColumn("type", TEXT, "")
	catalog.DefineColumn("id", KEY, "")
	catalog.DefineColumn("version", INT, 0)
	catalog.DefineColumn("definition", BYTES, []byte{})
	catalog.DefinePrimaryKeys("type", "id")
	catalog.IsCore = true
	if err = catalog.Init(); err != nil {
		return err
	}

	return nil
}

/**
* setCatalog
* @param tp, id string, version int, obj any
* @return error
**/
func setCatalog(tp, id string, version int, obj any) error {
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
			"id":         id,
			"version":    version,
			"definition": bt,
		}).
		BeforeInsertOrUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("created_at", now)
			new.Set("updated_at", now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("updated_at", now)
			return nil
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* getCatalog
* @param tp, id string, dest any
* @return bool, error
**/
func getCatalog(tp, id string, dest any) (bool, error) {
	item, err := catalog.
		Select().
		Where(Eq("type", tp)).
		And(Eq("id", id)).
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
* @param tp, id string
* @return error
**/
func deleteCatalog(tp, id string) error {
	_, err := catalog.
		Delete().
		Where(Eq("type", tp)).
		And(Eq("id", id)).
		One()
	if err != nil {
		return err
	}

	return nil
}
