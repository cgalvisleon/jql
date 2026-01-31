package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/et"
)

/**
* defaultChain
* @params et.Json
* @return string, error
**/
func defaultChain(params et.Json) (string, error) {
	username := params.Str("username")
	password := params.Str("password")
	host := params.Str("host")
	port := params.Int("port")
	app := params.Str("app")
	return fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, username, password, host, port, "postgres", app), nil
}

/**
* chain
* @params et.Json
* @return string, error
**/
func chain(params et.Json) (string, error) {
	username := params.Str("username")
	password := params.Str("password")
	host := params.Str("host")
	port := params.Int("port")
	app := params.Str("app")
	database := params.Str("database")
	if database == "" {
		return "", fmt.Errorf("database is required")
	}
	if host == "" {
		return "", fmt.Errorf("host is required")
	}
	if port == 0 {
		return "", fmt.Errorf("port is required")
	}
	if username == "" {
		return "", fmt.Errorf("username is required")
	}
	if password == "" {
		return "", fmt.Errorf("password is required")
	}
	if app == "" {
		return "", fmt.Errorf("app is required")
	}

	result := fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, username, password, host, port, database, app)
	return result, nil
}

/**
* connectTo
* @param chain string
* @return *sql.DB, error
**/
func connectTo(chain string) (*sql.DB, error) {
	db, err := sql.Open(driver, chain)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
