package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jql/jql"
)

/**
* buildModel
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildModel(model *jql.Model) (string, error) {
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

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildTriggerBeforeInsert(model)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	return sql, nil
}

/**
* buildSchema
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildSchema(model *jql.Model) (string, error) {
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
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildTable(model *jql.Model) (string, error) {
	getType := func(tp jql.TypeData) string {
		types := map[jql.TypeData]string{
			jql.ANY:      "VARCHAR(250)",
			jql.BYTES:    "BYTEA",
			jql.INT:      "BIGINT",
			jql.FLOAT:    "DOUBLE PRECISION",
			jql.KEY:      "VARCHAR(80)",
			jql.TEXT:     "VARCHAR(250)",
			jql.MEMO:     "TEXT",
			jql.JSON:     "JSONB",
			jql.DATETIME: "TIMESTAMP",
			jql.BOOLEAN:  "BOOLEAN",
			jql.GEOMETRY: "JSONB",
		}

		if t, ok := types[tp]; ok {
			return t
		}

		return "VARCHAR(250)"
	}

	defaultValue := func(tp jql.TypeData) string {
		values := map[jql.TypeData]string{
			jql.ANY:      "",
			jql.BYTES:    "''",
			jql.INT:      "0",
			jql.FLOAT:    "0.0",
			jql.KEY:      "''",
			jql.TEXT:     "''",
			jql.MEMO:     "''",
			jql.JSON:     "'{}'",
			jql.DATETIME: "NOW()",
			jql.BOOLEAN:  "FALSE",
			jql.GEOMETRY: "'{}'",
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
		if tpColumn != jql.COLUMN {
			continue
		}
		tpData := column.TypeData
		tp := getType(tpData)
		df := column.Default
		if df == nil {
			df = defaultValue(tpData)
		}
		def := fmt.Sprintf("\n\t%s %s DEFAULT %s", column.Name, tp, df)
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
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildPrimaryKeys(model *jql.Model) (string, error) {
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
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildForeignKeys(model *jql.Model) (string, error) {
	result := ""
	for name, foreignKey := range model.ForeignKeys {
		to := foreignKey.To.Table
		fks := ""
		ks := ""
		for k, fk := range foreignKey.Keys {
			fks = strs.Append(fks, fmt.Sprintf("%s", k), ", ")
			ks = strs.Append(ks, fmt.Sprintf("%s", fk), ", ")
		}
		def := fmt.Sprintf("ALTER TABLE IF EXISTS %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY(%s) REFERENCES %s(%s)", model.Table, model.Name, name, fks, to, ks)
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
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildIndexes(model *jql.Model) (string, error) {
	if len(model.Indexes) == 0 {
		return "", nil
	}

	table := model.Table
	name := model.Name
	result := ""
	for _, v := range model.Indexes {
		def := fmt.Sprintf("idx_%s_%s", name, v)
		if v == jql.SOURCE {
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
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildUniqueIndex(model *jql.Model) (string, error) {
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

/**
* buildTriggerBeforeInsert
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) buildTriggerBeforeInsert(model *jql.Model) (string, error) {
	if model.IndexField == nil {
		return "", nil
	}

	isCore := model.IsCore
	if isCore {
		return "", nil
	}

	table := model.Table
	result := fmt.Sprintf(`
	DROP TRIGGER IF EXISTS RECORDS_SET ON %s CASCADE;
	CREATE TRIGGER RECORDS_SET
	AFTER INSERT OR UPDATE OR DELETE ON %s
	FOR EACH ROW
	EXECUTE FUNCTION core.after_records();
	`, table, table)

	return result, nil
}

/**
* mutateModel
* @param model *jql.Model
* @return (string, error)
**/
func (s *Driver) mutateModel(model *jql.Model) (string, error) {
	return "", nil
}
