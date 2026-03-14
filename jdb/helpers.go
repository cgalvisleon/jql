package jdb

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

/**
* SQLParse
* @param sql string
* @param args ...any
* @return string
**/
func SQLParse(sql string, args ...any) string {
	for i := range args {
		old := fmt.Sprintf(`$%d`, i+1)
		new := fmt.Sprintf(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := fmt.Sprintf(`{$%d}`, i+1)
		new := fmt.Sprintf(`%v`, Quoted(arg))
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* SQLUnQuote
* @param sql string
* @param args ...any
* @return string
**/
func SQLUnQuote(sql string, args ...any) string {
	for i := range args {
		old := fmt.Sprintf(`$%d`, i+1)
		new := fmt.Sprintf(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := fmt.Sprintf(`{$%d}`, i+1)
		new := fmt.Sprintf(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* Quote
* @param val any
* @return any
**/
func Quoted(val any) any {
	format := `'%v'`
	switch v := val.(type) {
	case string:
		return fmt.Sprintf(format, v)
	case int:
		return v
	case float64:
		return v
	case float32:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case bool:
		return v
	case time.Time:
		return fmt.Sprintf(format, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return fmt.Sprintf(format, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(format, et.Json(v).ToString())
	case []string, []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote, type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return strs.Format(format, `[]`)
		}
		return fmt.Sprintf(format, string(bt))
	case []uint8:
		b := []byte(val.([]uint8))
		return fmt.Sprintf("'\\x%s'", hex.EncodeToString(b))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Quote, type:%v, value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* Literal
* @param val any
* @return any
**/
func Literal(val any) any {
	format := `"%v"`
	switch v := val.(type) {
	case string:
		return fmt.Sprintf(format, v)
	case int:
		return v
	case float64:
		return v
	case float32:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case bool:
		return v
	case time.Time:
		return fmt.Sprintf(format, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return fmt.Sprintf(format, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(format, et.Json(v).ToString())
	case []string, []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote, type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return strs.Format(format, `[]`)
		}
		return fmt.Sprintf(format, string(bt))
	case []uint8:
		b := []byte(val.([]uint8))
		return fmt.Sprintf("'\\x%s'", hex.EncodeToString(b))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Quote, type:%v, value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* insertBeforeLast
* @param s []T, v T
* @return []T
**/
func insertBeforeLast[T any](s []T, v T) []T {
	if len(s) < 1 {
		return append(s, v)
	}
	if len(s) == 1 {
		return []T{v, s[0]}
	}

	s = append(s[:len(s)-1], append([]T{v}, s[len(s)-1:]...)...)
	return s
}

/**
* RowsToItems
* @param rows *sql.Rows
* @return et.Items
**/
func RowsToItems(rows *sql.Rows) et.Items {
	defer rows.Close()

	result := et.Items{Result: []et.Json{}}
	append := func(item et.Json) {
		result.Add(item)
	}

	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		if len(item) == 1 {
			for _, v := range item {
				switch val := v.(type) {
				case et.Json:
					append(val)
				case map[string]interface{}:
					append(et.Json(val))
				default:
					append(item)
				}
			}
		} else {
			append(item)
		}
	}

	return result
}

/**
* findFieldByName
* @param froms []*From, name string // from.name:as|1:30
* @return *Field
**/
func findFieldByStr(froms []*From, name string) *Field {
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`) // from.name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)                // from.name
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)            // agg(field):as
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)                            // agg(field)
	pattern5 := regexp.MustCompile(`^(\d+)\|(\d+)$`)                                  // page:rows

	split := strings.Split(name, "|")
	if len(split) == 2 {
		name = split[0]
		limit := split[1]
		result := findFieldByStr(froms, name)
		if result != nil {
			if pattern5.MatchString(limit) {
				matches := pattern5.FindStringSubmatch(limit)
				if len(matches) == 3 {
					page, err := strconv.Atoi(matches[1])
					if err != nil {
						page = 0
					}
					rows, err := strconv.Atoi(matches[2])
					if err != nil {
						rows = 0
					}
					result.Page = page
					result.Rows = rows
				}
			}
		}

		return result
	}

	if pattern1.MatchString(name) {
		matches := pattern1.FindStringSubmatch(name)
		if len(matches) == 4 {
			from := matches[1]
			name = matches[2]
			as := matches[3]
			var result *Field
			for _, f := range froms {
				if f.As == from {
					result = f.findField(name)
				} else if f.Name == from {
					result = f.findField(name)
				}
				if result != nil {
					result.From = f
					result.As = as
					return result
				}
			}
		}
	} else if pattern2.MatchString(name) {
		matches := pattern2.FindStringSubmatch(name)
		if len(matches) == 3 {
			from := matches[1]
			name = matches[2]
			as := matches[2]
			var result *Field
			for _, f := range froms {
				if f.As == from {
					result = f.findField(name)
				} else if f.Name == from {
					result = f.findField(name)
				}
				if result != nil {
					result.From = f
					result.As = as
					return result
				}
			}
		}
	} else if pattern3.MatchString(name) {
		matches := pattern3.FindStringSubmatch(name)
		if len(matches) == 4 {
			agg := matches[1]
			name = matches[2]
			as := matches[3]
			if !slices.Contains(Aggs, agg) {
				return nil
			}
			result := findFieldByStr(froms, name)
			if result != nil {
				result.TypeColumn = AGG
				result.Field = &Agg{
					Agg:   agg,
					Field: name,
				}
				result.As = as
				return result
			}
		}
	} else if pattern4.MatchString(name) {
		matches := pattern4.FindStringSubmatch(name)
		if len(matches) == 3 {
			agg := matches[1]
			name = matches[2]
			as := agg
			if !slices.Contains(Aggs, agg) {
				return nil
			}
			result := findFieldByStr(froms, name)
			if result != nil {
				result.TypeColumn = AGG
				result.Field = &Agg{
					Agg:   agg,
					Field: name,
				}
				result.As = as
				return result
			}
		}
	} else {
		for _, f := range froms {
			result := f.findField(name)
			if result != nil {
				result.From = f
				return result
			}
		}
	}

	return nil
}
