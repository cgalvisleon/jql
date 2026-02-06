package jql

type Schema struct {
	Database string           `json:"-"`
	Name     string           `json:"name"`
	Models   map[string]*From `json:"models"`
	db       *DB              `json:"-"`
}
