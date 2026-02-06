package jql

import (
	"fmt"
	"net/http"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/utility"
)

var (
	ErrNotInserted = fmt.Errorf("record not inserted")
	ErrNotUpdated  = fmt.Errorf("record not updated")
	ErrNotFound    = fmt.Errorf("record not found")
	ErrNotUpserted = fmt.Errorf("record not inserted or updated")
	ErrDuplicate   = fmt.Errorf("record duplicate")
)

/**
* ConnectTo
* @param name string, params Connection
* @return (*DB, error)
**/
func ConnectTo(name string, params et.Json) (*DB, error) {
	return getDb(name, params)
}

/**
* LoadTo
* @param name string
* @return (*DB, error)
**/
func LoadTo(name string, host string, port int) (*DB, error) {
	params := et.Json{
		"driver":   envar.GetStr("DB_DRIVER", "postgres"),
		"database": name,
		"host":     host,
		"port":     port,
		"username": envar.GetStr("DB_USERNAME", "test"),
		"password": envar.GetStr("DB_PASSWORD", "test"),
		"app":      envar.GetStr("DB_APP", "jql"),
		"version":  envar.GetInt("DB_VERSION", 15),
	}

	return ConnectTo(name, params)
}

/**
* Load
* @return (*DB, error)
**/
func Load() (*DB, error) {
	name := envar.GetStr("DB_NAME", "josephine")
	host := envar.GetStr("DB_HOST", "localhost")
	port := envar.GetInt("DB_PORT", 5432)
	return LoadTo(name, host, port)
}

/**
* NewDb
* @param name string, host string, port int
* @return (*DB, error)
**/
func NewDb(name string, host string, port int) (*DB, error) {
	result, err := LoadTo(name, host, port)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* GetDb
* @param name string
* @return (*DB, error)
**/
func GetDb(name string) (*DB, error) {
	result, ok := dbs[name]
	if ok {
		return result, nil
	}

	exists, err := getCatalog("db", name, &result)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrDbNotFound
	}

	return result, nil
}

/**
* GetModel
* @param database, schema, name string
* @return (*Model, error)
**/
func GetModel(database, schema, name string) (*Model, error) {
	db, err := GetDb(database)
	if err != nil {
		return nil, err
	}

	result, err := db.GetModel(schema, name)
	if err != nil {
		return nil, err
	}

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

/**
* DeleteModel
* @param database, schema, name string
* @return error
**/
func DeleteModel(database, schema, name string) error {
	db, err := GetDb(database)
	if err != nil {
		return err
	}

	err = db.deleteModel(schema, name)
	if err != nil {
		return err
	}

	return nil
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func Define(definition et.Json) (*Model, error) {
	database := definition.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "database")
	}

	db, err := GetDb(database)
	if err != nil {
		return nil, err
	}

	return db.Define(definition)
}

/**
* Query
* @param query et.Json
* @return (et.Items, error)
**/
func Query(query et.Json) (et.Items, error) {
	database := query.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return et.Items{}, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED, "database")
	}

	db, err := GetDb(database)
	if err != nil {
		return et.Items{}, err
	}

	insert := query.Json("insert")
	if !insert.IsEmpty() {
		return db.Insert(insert)
	}

	update := query.Json("update")
	if !update.IsEmpty() {
		return db.Update(update)
	}

	delete := query.Json("delete")
	if !delete.IsEmpty() {
		return db.Delete(delete)
	}

	upsert := query.Json("upsert")
	if !upsert.IsEmpty() {
		return db.Upsert(upsert)
	}

	return db.Select(query)
}

/**
* HttpDefine
* @param w http.ResponseWriter, r *http.Request
* @return
**/
func HttpDefine(w http.ResponseWriter, r *http.Request) {
	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	result, err := Define(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result.ToJson())
}

/**
* HttpQuery
* @param w http.ResponseWriter, r *http.Request
* @return
**/
func HttpQuery(w http.ResponseWriter, r *http.Request) {
	body, err := response.GetBody(r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	result, err := Query(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEMS(w, r, http.StatusOK, result)
}

		return
	}

	response.ITEMS(w, r, http.StatusOK, result)
}
