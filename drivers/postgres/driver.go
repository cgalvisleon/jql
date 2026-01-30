package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jql/jql"
)

var driver = "postgres"

func init() {
	jql.Register(driver, new)
}

type Driver struct {
	name       string      `json:"-"`
	db         *jql.DB     `json:"-"`
	connection *Connection `json:"-"`
}

/**
* newDriver
* @param database *jql.DB
* @return jql.Driver
**/
func new() jql.Driver {
	result := &Driver{}
	return result
}

/**
* Connect
* @param connection jql.ConnectParams
* @return *sql.DB, error
**/
func (s *Driver) Connect(database *jql.DB) (*sql.DB, error) {
	s.database = database
	s.name = database.Name

	defaultChain, err := s.connection.defaultChain()
	if err != nil {
		return nil, err
	}

	db, err := connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	err = CreateDatabase(db, database.Name)
	if err != nil {
		return nil, err
	}

	if db != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
	}

	s.connection.Database = database.Name
	chain, err := s.connection.chain()
	if err != nil {
		return nil, err
	}

	db, err = connectTo(chain)
	if err != nil {
		return nil, err
	}

	if database.UseCore {
		if err := TriggerRecords(db); err != nil {
			return nil, err
		}
	}

	logs.Logf(driver, `Connected to %s:%s:%d`, s.connection.Host, s.connection.Database, s.connection.Port)

	return db, nil
}

/**
* Load
* @param model *Model
* @return (string, error)
**/
func (s *Driver) Load(model *jql.Model) (string, error) {
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
func (s *Driver) Mutate(model *jql.Model) (string, error) {
	return "", nil
}

/**
* Query
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) Query(ql *jql.Ql) (string, error) {
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
* @param command *jql.Cmd
* @return (string, error)
**/
func (s *Driver) Command(cmd *jql.Cmd) (string, error) {
	result, err := s.buildCommand(cmd)
	if err != nil {
		return "", err
	}

	if cmd.IsDebug {
		logs.Debug("command:\n", result)
	}

	return result, nil
}
