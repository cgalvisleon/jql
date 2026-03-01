package jql

import (
	"net/http"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/jql/jdb"
)

/**
* NewModel
* @param db *DB, schema, name string, version int
* @return (*Model, error)
**/
func NewModel(db *DB, schema, name string, version int) (*Model, error) {
	result, err := db.Define(et.Json{
		"schema":  schema,
		"name":    name,
		"version": version,
	})

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
* SetSerie
* @param tag, format string, value int
* @return error
**/
func SetSerie(tag, format string, value int) error {
	return jdb.SetSerie(tag, format, value)
}

/**
* GetSerie
* @param tag string
* @return (int, string, error)
**/
func GetSerie(tag string) (int, string, error) {
	return jdb.GetSerie(tag)
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
