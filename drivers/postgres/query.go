package postgres

import (
	"fmt"
	"strings"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jql/jql"
)

/**
* Query
* @param query *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildQuery(ql *jql.Ql) (string, error) {
	if ql.IsDebug {
		logs.Debug("query:", ql.ToJson().ToString())
	}

	sql, err := s.buildSelect(ql)
	if err != nil {
		return "", err
	}

	sql = fmt.Sprintf("SELECT %s", sql)
	def, err := s.buildFrom(ql)
	if err != nil {
		return "", err
	}

	def = fmt.Sprintf("FROM %s", def)
	sql = strs.Append(sql, def, "\n")
	def, err = s.buildJoins(ql)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("JOIN %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	wheres := ql.Wheres.Conditions
	if len(wheres) > 0 {
		def, err = s.buildWhere(wheres)
		if err != nil {
			return "", err
		}

		if def != "" {
			def = fmt.Sprintf("WHERE %s", def)
			sql = strs.Append(sql, def, "\n")
		}
	}

	def, err = s.buildGroupBy(ql)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("GROUP BY %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildWhere(ql.Havings.Conditions)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("HAVING %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildOrderBy(ql)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("ORDER BY %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildLimit(ql)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	if ql.Type == jql.EXISTS {
		return fmt.Sprintf("SELECT EXISTS(%s);", sql), nil
	} else {
		return fmt.Sprintf("%s;", sql), nil
	}
}

/**
* buildSelect
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildSelect(ql *jql.Ql) (string, error) {
	if ql.Type == jql.EXISTS {
		return "", nil
	}

	if ql.Type == jql.COUNTED {
		return "COUNT(*) AS count", nil
	}

	result := ""
	if ql.Type == jql.DATA {
		if len(ql.Selects) == 0 {
			hiddens := ql.Hidden
			hiddens = append(hiddens, jql.SOURCE)
			def := fmt.Sprintf("to_jsonb(A) - ARRAY[%s]", strings.Join(hiddens, ", "))
			result = strs.Append(result, def, "||")
		} else {
			selects := map[string]string{}
			atribs := map[string]string{}
			for _, fld := range ql.Selects {
				if fld.TypeColumn == jql.COLUMN {
					as := FieldAs(fld)
					selects[fld.As] = as
				} else if fld.TypeColumn == jql.AGG {
					as := FieldAs(fld)
					selects[fld.As] = as
				} else if fld.TypeColumn == jql.ATTRIB {
					as := FieldAs(fld)
					atribs[fld.As] = as
				}
			}

			if len(atribs) == 0 {
				result = fmt.Sprintf("\n%s", jql.SOURCE)
			} else {
				for k, v := range atribs {
					def := fmt.Sprintf("\n'%s', %s", k, v)
					result = strs.Append(result, def, ", ")
				}

				if result != "" {
					result = fmt.Sprintf("\n\tjsonb_build_object(%s\n)", result)
				}
			}

			sel := ""
			for k, v := range selects {
				def := fmt.Sprintf("\n'%s',  %s", k, v)
				if v == "" {
					def = fmt.Sprintf("\n'%s',  %s", k, k)
				}
				sel = strs.Append(sel, def, ", ")
			}

			if sel != "" {
				result = fmt.Sprintf("%s||jsonb_build_object(%s\n)", result, sel)
			}
		}

		return fmt.Sprintf("%s AS result", result), nil
	}

	if len(ql.Selects) == 0 {
		hiddens := ql.Hidden
		if len(hiddens) > 0 {
			result += fmt.Sprintf("to_jsonb(A) - ARRAY[%s]", strings.Join(hiddens, ", "))
		} else {
			result += "A.*"
		}
	} else {
		selects := map[string]string{}
		for _, fld := range ql.Selects {
			if fld.TypeColumn == jql.COLUMN {
				as := FieldAs(fld)
				selects[fld.As] = as
			}
		}
		for k, v := range selects {
			def := fmt.Sprintf("\n%s AS %s", v, k)
			if k == v {
				def = fmt.Sprintf("\n%s", v)
			} else if v == "" {
				def = fmt.Sprintf("\n%s", v)
			}
			result = strs.Append(result, def, ", ")
		}
	}

	return result, nil
}

/**
* buildFrom
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildFrom(ql *jql.Ql) (string, error) {
	result := ""

	if len(ql.Froms) == 0 {
		return result, fmt.Errorf(jql.MSG_FROM_REQUIRED)
	}

	for _, from := range ql.Froms {
		as := from.As
		table := from.Table
		def := fmt.Sprintf("%s AS %s", table, as)
		if as == table {
			def = fmt.Sprintf("%s", table)
		}

		result = strs.Append(result, def, ", ")
		break
	}

	return result, nil
}

/**
* buildJoins
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildJoins(ql *jql.Ql) (string, error) {
	result := ""

	if len(ql.Joins) == 0 {
		return result, nil
	}

	for _, join := range ql.Joins {
		def := ""
		for k, v := range join.Keys {
			if len(def) == 0 {
				def = fmt.Sprintf("%s AS %s ON %s = %s", join.To.As, join.To.As, k, v)
			} else {
				def = fmt.Sprintf("%s AND %s = %s", def, k, v)
			}
		}

		if join.Type == jql.LEFT {
			result = strs.Append(result, def, "\nLEFT JOIN ")
		} else if join.Type == jql.RIGHT {
			result = strs.Append(result, def, "\nRIGHT JOIN ")
		} else if join.Type == jql.FULL {
			result = strs.Append(result, def, "\nFULL JOIN ")
		} else {
			result = strs.Append(result, def, "\nJOIN ")
		}
	}

	return fmt.Sprintf("%s", result), nil
}

/**
* buildCondition
* @param cond *jql.Condition
* @return string
**/
func (s *Driver) buildCondition(cond *jql.Condition) string {
	key := FieldAs(cond.Field)
	value := jql.Quoted(cond.Value)
	switch cond.Operator {
	case jql.OpEq:
		return fmt.Sprintf("%s = %v", key, value)
	case jql.OpNeg:
		return fmt.Sprintf("%s != %v", key, value)
	case jql.OpLess:
		return fmt.Sprintf("%s < %v", key, value)
	case jql.OpLessEq:
		return fmt.Sprintf("%s <= %v", key, value)
	case jql.OpMore:
		return fmt.Sprintf("%s > %v", key, value)
	case jql.OpMoreEq:
		return fmt.Sprintf("%s >= %v", key, value)
	case jql.OpLike:
		return fmt.Sprintf("%s LIKE %v", key, value)
	case jql.OpIn:
		return fmt.Sprintf("%s IN %v", key, value)
	case jql.OpNotIn:
		return fmt.Sprintf("%s NOT IN %v", key, value)
	case jql.OpIs:
		return fmt.Sprintf("%s IS %v", key, value)
	case jql.OpIsNot:
		return fmt.Sprintf("%s IS NOT %v", key, value)
	case jql.OpNull:
		return fmt.Sprintf("%s IS NULL", key)
	case jql.OpNotNull:
		return fmt.Sprintf("%s IS NOT NULL", key)
	case jql.OpBetween:
		vals := cond.Value.([]interface{})
		return fmt.Sprintf("%s BETWEEN %v AND %v", key, jql.Quoted(vals[0]), jql.Quoted(vals[1]))
	case jql.OpNotBetween:
		vals := cond.Value.([]interface{})
		return fmt.Sprintf("%s NOT BETWEEN %v AND %v", key, jql.Quoted(vals[0]), jql.Quoted(vals[1]))
	}

	return ""
}

/**
* buildWhere
* @param wheres []jql.Condition
* @return (string, error)
**/
func (s *Driver) buildWhere(wheres []*jql.Condition) (string, error) {
	result := ""

	for i, cond := range wheres {
		if i == 0 {
			result = s.buildCondition(cond)
		} else if cond.Connector == jql.OR {
			result = fmt.Sprintf("%s\nOR %s", result, s.buildCondition(cond))
		} else {
			result = fmt.Sprintf("%s\nAND %s", result, s.buildCondition(cond))
		}
	}

	return result, nil
}

/**
* buildGroupBy
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildGroupBy(ql *jql.Ql) (string, error) {
	result := ""

	if len(ql.GroupsBy) == 0 {
		return result, nil
	}

	for _, v := range ql.GroupsBy {
		as := FieldAs(v)
		def := fmt.Sprintf("%s", as)
		result = strs.Append(result, def, ", ")
	}

	return result, nil
}

/**
* buildOrderBy
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildOrderBy(ql *jql.Ql) (string, error) {
	asc := ""
	desc := ""
	for _, order := range ql.OrdersBy {
		as := FieldAs(order.Field)
		if order.Asc {
			asc = strs.Append(asc, as, ", ")
		} else {
			desc = strs.Append(desc, as, ", ")
		}
	}

	result := ""
	if asc != "" {
		result = fmt.Sprintf(`%s ASC`, asc)
	}

	if desc != "" {
		result = fmt.Sprintf(`%s DESC`, desc)
	}

	return result, nil
}

/**
* buildLimit
* @param ql *jql.Ql
* @return (string, error)
**/
func (s *Driver) buildLimit(ql *jql.Ql) (string, error) {
	result := ""

	if ql.Rows > ql.MaxRows {
		ql.Rows = ql.MaxRows
	} else if ql.Rows == 0 {
		ql.Rows = 1
	}

	if ql.Page == 0 {
		return fmt.Sprintf("LIMIT %d", ql.Rows), nil
	}

	offset := (ql.Page - 1) * ql.Rows
	result = fmt.Sprintf("%d OFFSET %d", ql.Rows, offset)
	return result, nil
}
