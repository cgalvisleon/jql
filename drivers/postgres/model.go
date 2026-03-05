package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jql/jdb"
)

/**
* buildModel
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildModel(model *jdb.Model) (string, error) {
	if model.IsDebug {
		logs.Debug("model:", model.ToJson().ToString())
	}

	exists, err := ExistTable(model.Db(), model.Schema, model.Name)
	if err != nil {
		return "", err
	}

	if exists {
		return "", nil
	}

	sql, err := s.buildSchema(model)
	if err != nil {
		return "", err
	}

	def, err := s.buildTable(model)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildIndexes(model)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildUniqueIndex(model)
	if err != nil {
		return "", err
	}

	def, err = s.buildForeignKeys(model)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	return sql, nil
}

/**
* buildSchema
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildSchema(model *jdb.Model) (string, error) {
	if !utility.ValidStr(model.Schema, 0, []string{}) {
		return "", fmt.Errorf(MSG_ATRIB_REQUIRED, "schema")
	}

	exist, err := ExistSchema(model.Db(), model.Schema)
	if err != nil {
		return "", err
	}

	if exist {
		return "", nil
	}

	return fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", model.Schema), nil
}

/**
* buildTable
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildTable(model *jdb.Model) (string, error) {
	getType := func(tp jdb.TypeData) string {
		types := map[jdb.TypeData]string{
			jdb.ANY:      "VARCHAR(250)",
			jdb.BYTES:    "BYTEA",
			jdb.INT:      "BIGINT",
			jdb.FLOAT:    "DOUBLE PRECISION",
			jdb.KEY:      "VARCHAR(80)",
			jdb.TEXT:     "VARCHAR(250)",
			jdb.MEMO:     "TEXT",
			jdb.JSON:     "JSONB",
			jdb.DATETIME: "TIMESTAMP",
			jdb.BOOLEAN:  "BOOLEAN",
			jdb.GEOMETRY: "JSONB",
		}

		if t, ok := types[tp]; ok {
			return t
		}

		return "VARCHAR(250)"
	}

	defaultValue := func(tp jdb.TypeData) string {
		values := map[jdb.TypeData]string{
			jdb.ANY:      "",
			jdb.BYTES:    "''",
			jdb.INT:      "0",
			jdb.FLOAT:    "0.0",
			jdb.KEY:      "''",
			jdb.TEXT:     "''",
			jdb.MEMO:     "''",
			jdb.JSON:     "'{}'",
			jdb.DATETIME: "NOW()",
			jdb.BOOLEAN:  "FALSE",
			jdb.GEOMETRY: "'{}'",
		}

		if t, ok := values[tp]; ok {
			return t
		}

		return ""
	}

	columns := model.Columns
	columnsDef := ""
	for _, column := range columns {
		tpColumn := column.TypeColumn
		if tpColumn != jdb.COLUMN {
			continue
		}
		tpData := column.TypeData
		tp := getType(tpData)
		df := column.Default
		switch v := df.(type) {
		case string:
			if v == "" {
				df = defaultValue(tpData)
			}
		default:
			df = defaultValue(tpData)
		}
		def := fmt.Sprintf("\n\t%s %s DEFAULT %v", column.Name, tp, df)
		columnsDef = strs.Append(columnsDef, def, ",")
	}

	def, err := s.buildPrimaryKeys(model)
	if err != nil {
		return "", err
	}

	if def != "" {
		columnsDef = strs.Append(columnsDef, def, ",\n\t")
	}

	table := model.Table
	result := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table, columnsDef)

	return result, nil
}

/**
* buildPrimaryKeys
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildPrimaryKeys(model *jdb.Model) (string, error) {
	if len(model.PrimaryKeys) == 0 {
		return "", nil
	}

	columns := ""
	for _, v := range model.PrimaryKeys {
		columns = strs.Append(columns, v, ", ")
	}

	name := model.Name
	result := fmt.Sprintf("CONSTRAINT pk_%s PRIMARY KEY (%s)", name, columns)

	return result, nil
}

/**
* buildForeignKeys
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildForeignKeys(model *jdb.Model) (string, error) {
	result := ""
	for _, foreignKey := range model.ForeignKeys {
		name := fmt.Sprintf("fk_%s_%s", model.Name, foreignKey.To.Name)
		to := foreignKey.To.Table
		fks := ""
		ks := ""
		for k, fk := range foreignKey.Keys {
			fks = strs.Append(fks, fmt.Sprintf("%s", k), ", ")
			ks = strs.Append(ks, fmt.Sprintf("%s", fk), ", ")
		}
		def := fmt.Sprintf("ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY(%s) REFERENCES %s(%s)", model.Table, name, fks, to, ks)
		onDelete := foreignKey.OnDeleteCascade
		onUpdate := foreignKey.OnUpdateCascade
		if onDelete {
			def = strs.Append(def, fmt.Sprintf("ON DELETE CASCADE"), " ")
		}
		if onUpdate {
			def = strs.Append(def, fmt.Sprintf("ON UPDATE CASCADE"), " ")
		}
		def += ";"
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}

/**
* buildIndexes
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildIndexes(model *jdb.Model) (string, error) {
	if len(model.Indexes) == 0 {
		return "", nil
	}

	table := model.Table
	name := model.Name
	result := ""
	for _, v := range model.Indexes {
		def := fmt.Sprintf("idx_%s_%s", name, v)
		if v == jdb.SOURCE {
			def = fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s USING GIN (%s);", def, table, v)
		} else {
			def = fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(%s);", def, table, v)
		}
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}

/**
* buildUniqueIndex
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildUniqueIndex(model *jdb.Model) (string, error) {
	if len(model.Unique) == 0 {
		return "", nil
	}

	table := model.Table
	name := model.Name
	result := ""
	for _, v := range model.Unique {
		def := fmt.Sprintf("idx_%s_%s", name, v)
		def = fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s(%s);", def, table, v)
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}
