package jql

import (
	"fmt"
	"net/http"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jql/jdb"
)

var (
	ErrNotInserted = fmt.Errorf("record not inserted")
	ErrNotFound    = fmt.Errorf("record not found")
	ErrNotUpserted = fmt.Errorf("record not inserted or updated")
	ErrDuplicate   = fmt.Errorf("record duplicate")
)

/**
* ConnectTo
* @param name string, params Connection
* @return (*DB, error)
**/
func ConnectTo(name string, params et.Json) (*jdb.DB, error) {
	return jdb.LoadDb(name, params)
}

/**
* LoadTo
* @param name string
* @return *jdb.DB, error
**/
func LoadTo(name string, host string, port int) (*jdb.DB, error) {
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
* @return (*jdb.DB, error)
**/
func Load() (*jdb.DB, error) {
	name := envar.GetStr("DB_NAME", "josephine")
	host := envar.GetStr("DB_HOST", "localhost")
	port := envar.GetInt("DB_PORT", 5432)
	return LoadTo(name, host, port)
}

/**
* GetModel
* @param database, schema, name string
* @return *jdb.Model, error
**/
func GetModel(database, schema, name string) (*jdb.Model, error) {
	db, err := jdb.GetDb(database)
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

/**
* Define
* @param params et.Json
* @return *jdb.Model, error
**/
func Define(params et.Json) (*jdb.Model, error) {
	database := params.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(jdb.MSG_ATTRIBUTE_REQUIRED, "database")
	}

	db, err := jdb.GetDb(database)
	if err != nil {
		return nil, err
	}

	return db.Define(params)
}

/**
* Insert
* @param params et.Json
* @return (et.Items, error)
**/
func Insert(params et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* Update
* @param params et.Json
* @return (et.Items, error)
**/
func Update(params et.Json) (et.Items, error) {
	return et.Items{}, nil
}

func Delete(params et.Json) (et.Items, error) {
	return et.Items{}, nil
}

func Upsert(params et.Json) (et.Items, error) {
	return et.Items{}, nil
}

/**
* From
* @param query et.Json
* @return (et.Items, error)
*
 */
func From(query et.Json) (et.Items, error) {
	database := query.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return et.Items{}, fmt.Errorf(jdb.MSG_ATTRIBUTE_REQUIRED, "database")
	}

	db, err := jdb.GetDb(database)
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
