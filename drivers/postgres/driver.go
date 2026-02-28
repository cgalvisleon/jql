package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jql/jdb"
)

var driver = "postgres"

func init() {
	jdb.Register(driver, new)
}

type Driver struct{}

/**
* newDriver
* @param database *jdb.DB
* @return jdb.Driver
**/
func new() jdb.Driver {
	result := &Driver{}
	return result
}

/**
* Connect
* @param connection jdb.ConnectParams
* @return *sql.DB, error
**/
func (s *Driver) Connect(db *jdb.DB) (*sql.DB, error) {
	params := db.Params
	defaultChain, err := defaultChain(params)
	if err != nil {
		return nil, err
	}

	result, err := connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	err = CreateDatabase(result, db.Name)
	if err != nil {
		return nil, err
	}

	if result != nil {
		err := result.Close()
		if err != nil {
			return nil, err
		}
	}

	chain, err := chain(params)
	if err != nil {
		return nil, err
	}

	result, err = connectTo(chain)
	if err != nil {
		return nil, err
	}

	host := params.Str("host")
	port := params.Int("port")
	logs.Logf(driver, `Connected to %s:%s:%d`, host, db.Name, port)

	return result, nil
}

/**
* Load
* @param model *Model
* @return (string, error)
**/
func (s *Driver) Load(model *jdb.Model) (string, error) {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)
	result, err := s.buildModel(model)
	if err != nil {
		return "", err
	}

	if model.IsDebug {
		logs.Debug("model:\n", result)
	}

	logs.Logf(driver, MSG_LOAD_MODEL, model.Name, model.Version)
	return result, nil
}

/**
* Mutate
* @param model *Model
* @return (string, error)
**/
func (s *Driver) Mutate(model *jdb.Model) (string, error) {
	return "", nil
}

/**
* Query
* @param ql *jdb.Ql
* @return (string, error)
**/
func (s *Driver) Query(ql *jdb.Ql) (string, error) {
	result, err := s.buildQuery(ql)
	if err != nil {
		return "", err
	}

	if ql.IsDebug {
		logs.Debug("query:\n", result)
	}

	return result, nil
}

/**
* Cmd
* @param command *jdb.Cmd
* @return (string, error)
**/
func (s *Driver) Command(cmd *jdb.Cmd) (string, error) {
	result, err := s.buildCommand(cmd)
	if err != nil {
		return "", err
	}

	if cmd.IsDebug {
		logs.Debug("command:\n", result)
	}

	return result, nil
}
