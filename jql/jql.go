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
* @param name string, data et.Json
* @return (jdb.Cmd, error)
**/
func Insert(name string, data et.Json) jdb.Cmd {
	return jdb.Cmd{}
}

/**
* Update
* @param name string, data et.Json
* @return (jdb.Cmd, error)
**/
func Update(name string, data et.Json) jdb.Cmd {
	return jdb.Cmd{}
}

/**
* Delete
* @param name string
* @return jdb.Cmd
**/
func Delete(name string) jdb.Cmd {
	return jdb.Cmd{}
}

/**
* Upsert
* @param name string, data et.Json
* @return jdb.Cmd
**/
func Upsert(name string, data et.Json) jdb.Cmd {
	return jdb.Cmd{}
}

/**
* From
* @param name, as string
* @return jdb.Ql
*
 */
func From(name, as string) jdb.Ql {
	return jdb.Ql{}
}

/**
* Query
* @param params et.Json
* @return (et.Items, error)
**/
func Query(params et.Json) (et.Items, error) {
	return et.Items{}, nil
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
