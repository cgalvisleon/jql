package jql

import (
	"encoding/json"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
)

type TypeQuery string

const (
	SELECT  TypeQuery = "select"
	DATA    TypeQuery = "data"
	EXISTS  TypeQuery = "exists"
	COUNTED TypeQuery = "count"
)

type Orders struct {
	Field *Field `json:"field"`
	Asc   bool   `json:"asc"`
}

type Ql struct {
	DB       *DB                    `json:"-"`
	Type     TypeQuery              `json:"type"`
	Froms    []*From                `json:"froms"`
	Selects  []interface{}          `json:"select"`
	Hidden   []*Field               `json:"hidden"`
	Wheres   *Wheres                `json:"wheres"`
	Joins    []*Joins               `json:"joins"`
	Details  map[string]*Detail     `json:"details"`
	Rollups  map[string]*Detail     `json:"rollups"`
	Calcs    map[string]DataContext `json:"calcs"`
	GroupsBy []*Field               `json:"group_by"`
	Havings  *Wheres                `json:"having"`
	OrdersBy []*Orders              `json:"order_by"`
	Page     int                    `json:"page"`
	Rows     int                    `json:"rows"`
	MaxRows  int                    `json:"max_rows"`
	IsDebug  bool                   `json:"is_debug"`
	tx       *Tx                    `json:"-"`
}

/**
* newQuery
* @param model *Model, as string, tp TypeQuery
* @return *Ql
**/
func newQuery(model *Model, as string) *Ql {
	tp := SELECT
	if model.SourceField != "" {
		tp = DATA
	}
	from := model.from()
	from.As = as
	maxRows := envar.GetInt("MAX_ROWS", 100)
	result := &Ql{
		Type:     tp,
		DB:       model.db,
		Froms:    []*From{from},
		Joins:    make([]*Joins, 0),
		Selects:  make([]interface{}, 0),
		Hidden:   make([]*Field, 0),
		Details:  make(map[string]*Detail),
		Rollups:  make(map[string]*Detail),
		Calcs:    make(map[string]DataContext),
		GroupsBy: make([]*Field, 0),
		OrdersBy: make([]*Orders, 0),
		Page:     0,
		Rows:     0,
		MaxRows:  maxRows,
	}
	result.Wheres = newWhere()
	result.Wheres.setModel(model)
	result.Havings = newWhere()

	return result
}

/**
* serialize
* @return []byte, error
**/
func (s *Ql) serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Ql) ToJson() et.Json {
	bt, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* setDebug
* @param isDebug bool
* @return *Ql
**/
func (s *Ql) setDebug(isDebug bool) *Ql {
	s.IsDebug = isDebug
	return s
}

/**
* Debug
* @return *Ql
**/
func (s *Ql) Debug() *Ql {
	s.setDebug(true)
	return s
}

/**
* FindField
* @param froms []*From, name string // from.name:as|1:30
* @return *Field
**/
func findField(froms []*From, name string) *Field {
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`) // from.name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)                // from.name
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)            // agg(field):as
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)                            // agg(field)
	pattern5 := regexp.MustCompile(`^(\d+)\|(\d+)$`)                                  // page:rows

	split := strings.Split(name, "|")
	if len(split) == 2 {
		name = split[0]
		limit := split[1]
		result := findField(froms, name)
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
			result := findField(froms, name)
			if result != nil {
				result.TypeColumn = TpAggregation
				result.Aggregation = GetAggregation(aggregation)
				result.As = as
				return result
			}
		}
	} else if pattern4.MatchString(name) {
		matches := pattern4.FindStringSubmatch(name)
		if len(matches) == 3 {
			aggregation := matches[1]
			name = matches[2]
			result := FindField(froms, name)
			if result != nil {
				result.TypeColumn = TpAggregation
				result.Aggregation = GetAggregation(aggregation)
				result.As = aggregation
				return result
			}
		}
	} else {
		if len(froms) == 0 {
			return nil
		}
		from := froms[0]
		result := from.FindField(name)
		if result != nil {
			return result
		}
	}

	return nil
}

/**
* Select
* @return *Ql
**/
func (s *Ql) Select(fields ...interface{}) *Ql {
	if len(s.Froms) == 0 {
		return s
	}

	if len(fields) == 0 {
		s.Selects = make([]interface{}, 0)
		for _, from := range s.Froms {
			for _, field := range from.Fields {
				if field.TypeColumn == COLUMN {
					field.From = from
					s.Selects = append(s.Selects, field)
				}
			}
		}
		return s
	}

	for _, fld := range fields {
		switch v := fld.(type) {
		case string:
			f := findField(s.Froms, v)
			f.From = s.Froms[0]
			s.Selects = append(s.Selects, f)
		case *Field:
			v.From = s.Froms[0]
			s.Selects = append(s.Selects, v)
		case *Agg:
			ql.Selects = append(ql.Selects, v)
		}
	}
	return ql
}

/**
* Join
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) join(tp TypeJoin, model *Model, as string, keys map[string]string) *Ql {
	from := newFrom(model, as)
	s.Froms = append(s.Froms, from)

	rKeys := make(map[string]string)
	for k, fk := range keys {
		field := FindField(s.Froms, k)
		if field != nil {
			k = field.AS()
		}

		field = FindField(s.Froms, fk)
		if field != nil {
			fk = field.AS()
		}

		rKeys[k] = fk
	}

	join := newJoins(tp, from, rKeys)
	s.Joins = append(s.Joins, join)

	return s
}

/**
* Join
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) Join(model *Model, as string, keys map[string]string) *Ql {
	return s.join(TpJoin, model, as, keys)
}

/**
* LeftJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) LeftJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(TpLeft, model, as, keys)
}

/**
* RightJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) RightJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(TpRight, model, as, keys)
}

/**
* FullJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) FullJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(TpFull, model, as, keys)
}

/**
* Where
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) Where(condition *Condition) *Ql {
	s.Wheres.Add(condition)
	return s
}

/**
* WhereByKeys
* @param keys et.Json
* @return *Ql
**/
func (s *Ql) WhereByKeys(keys et.Json) *Ql {
	for k, v := range keys {
		s.Wheres.Add(Eq(k, v))
	}
	return s
}

/**
* WhereByConditions
* @param conditions []*Condition
* @return *Ql
**/
func (s *Ql) WhereByConditions(conditions []*Condition) *Ql {
	for _, condition := range conditions {
		s.Wheres.Add(condition)
	}
	return s
}

/**
* WhereByJson
* @param jsons []et.Json
* @return *Ql
**/
func (s *Ql) WhereByJson(jsons []et.Json) *Ql {
	s.Wheres.ByJson(jsons)
	return s
}

/**
* And
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) And(condition *Condition) *Ql {
	s.Wheres.Add(condition)
	return s
}

/**
* Or
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) Or(condition *Condition) *Ql {
	s.Wheres.Add(condition)
	return s
}

/**
* GroupsBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) GroupBy(fields ...string) *Ql {
	for _, name := range fields {
		fld := FindField(s.Froms, name)
		if fld != nil {
			s.GroupsBy = append(s.GroupsBy, fld)
		}
	}
	return s
}

/**
* Having
* @param condition []*Condition
* @return *Ql
**/
func (s *Ql) Having(condition []*Condition) *Ql {
	for _, cnd := range condition {
		s.Havings.Add(cnd)
	}
	return s
}

/**
* HavingsByJson
* @param jsons []et.Json
* @return *Ql
**/
func (s *Ql) HavingsByJson(jsons []et.Json) *Ql {
	s.Havings.ByJson(jsons)
	return s
}

/**
* ordersBy
* @param asc bool, fields ...string
* @return *Ql
**/
func (s *Ql) ordersBy(asc bool, fields ...string) *Ql {
	for _, name := range fields {
		fld := FindField(s.Froms, name)
		if fld != nil {
			s.OrdersBy = append(s.OrdersBy, &Orders{Field: fld, Asc: asc})
		}
	}
	return s
}

/**
* OrderBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderBy(fields ...string) *Ql {
	return s.ordersBy(true, fields...)
}

/**
* OrderByAsc
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderByAsc(fields ...string) *Ql {
	return s.ordersBy(true, fields...)
}

/**
* OrderByDesc
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderByDesc(fields ...string) *Ql {
	return s.ordersBy(false, fields...)
}

/**
* Hiddens
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Hiddens(fields ...string) *Ql {
	for _, name := range fields {
		fld := FindField(s.Froms, name)
		if fld != nil {
			s.Hidden = append(s.Hidden, fld)
		}
	}
	return s
}

/**
* prepare
**/
func (s *Ql) prepare() {
	if len(s.Selects) == 0 {
		return
	}

	for _, hidden := range s.Hidden {
		idx := slices.IndexFunc(s.Selects, func(fld *Field) bool { return fld.AS() == hidden.AS() })
		if idx != -1 {
			s.Selects = append(s.Selects[:idx], s.Selects[idx+1:]...)
		}
	}
}

/**
* AllTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Ql) AllTx(tx *Tx) (et.Items, error) {
	s.prepare()
	return s.DB.Query(s)
}

/**
* All
* @return et.Items, error
**/
func (s *Ql) All() (et.Items, error) {
	return s.AllTx(nil)
}

/**
* LimitTx
* @param tx *Tx, page, rows int
* @return *Ql
**/
func (s *Ql) LimitTx(tx *Tx, page, rows int) (et.Items, error) {
	s.Page = page
	s.Rows = rows
	return s.AllTx(tx)
}

/**
* Limit
* @param page, rows int
* @return *Ql
**/
func (s *Ql) Limit(page, rows int) (et.Items, error) {
	return s.LimitTx(nil, page, rows)
}

/**
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Ql) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.AllTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* One
* @param tx *Tx
* @return et.Item, error
**/
func (s *Ql) One() (et.Item, error) {
	return s.OneTx(nil)
}

/**
* ItExistsTx
* @param tx *Tx
* @return bool, error
**/
func (s *Ql) ItExistsTx(tx *Tx) (bool, error) {
	s.Type = TpExists
	result, err := s.AllTx(tx)
	if err != nil {
		return false, err
	}

	exists := result.First().Bool("exists")
	return exists, nil
}

/**
* ItExists
* @return bool, error
**/
func (s *Ql) ItExists() (bool, error) {
	return s.ItExistsTx(nil)
}

/**
* CountTx
* @param tx *Tx
* @return int, error
**/
func (s *Ql) CountTx(tx *Tx) (int, error) {
	s.Type = TpCounted
	result, err := s.AllTx(tx)
	if err != nil {
		return 0, err
	}

	count := result.First().Int("count")
	return count, nil
}

/**
* Count
* @return int, error
**/
func (s *Ql) Count() (int, error) {
	return s.CountTx(nil)
}
