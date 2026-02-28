package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var series *Model

/**
* defineSeries
* @param db *DB
* @return error
**/
func defineSeries(db *DB) error {
	if series != nil {
		return nil
	}

	var err error
	series, err = db.NewModel("core", "series", 1)
	if err != nil {
		return err
	}
	series.defineCreatedAtField()
	series.defineUpdatedAtField()
	series.DefineColumn("tag", TEXT, "")
	series.DefineColumn("format", TEXT, "")
	series.DefineColumn("value", INT, 0)
	series.DefinePrimaryKeys("tag")
	series.IsCore = true
	if err = series.Init(); err != nil {
		return err
	}

	return nil
}

/**
* SetSerie
* @param tag, format string, value int
* @return error
**/
func SetSerie(tag, format string, value int) error {
	if series == nil {
		return nil
	}

	now := timezone.Now()
	_, err := series.
		Upsert(et.Json{
			"tag":    tag,
			"format": format,
			"value":  value,
		}).
		BeforeInsert(func(tx *Tx, old, new et.Json) error {
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
* GetSerie
* @param tag string
* @return et.Json, error
**/
func GetSerie(tag string) (et.Json, error) {
	if series == nil {
		return et.Json{}, nil
	}

	now := timezone.Now()
	item, err := series.
		Upsert(et.Json{
			"tag": tag,
		}).
		BeforeInsert(func(tx *Tx, old, new et.Json) error {
			new.Set("created_at", now)
			new.Set("updated_at", now)
			new.Set("format", "%08d")
			new.Set("value", 1)

			return nil
		}).
		BeforeUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("updated_at", now)
			new.Set("value", old.Int("value")+1)
			return nil
		}).
		Where(Eq("tag", tag)).
		One()
	if err != nil {
		return et.Json{}, err
	}

	return et.Json{
		"value":  item.Int("value"),
		"format": item.String("format"),
	}, nil
}

/**
* DeleteSerie
* @param tag string
* @return error
**/
func DeleteSerie(tag string) error {
	_, err := series.
		Delete().
		Where(Eq("tag", tag)).
		One()
	if err != nil {
		return err
	}

	return nil
}
