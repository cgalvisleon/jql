package jdb

import (
	"encoding/json"

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
	Type     TypeQuery              `json:"type"`
	Froms    []*From                `json:"froms"`
	Selects  []*Field               `json:"select"`
	Hiddens  []string               `json:"hidden"`
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
	db       *DB                    `json:"-"`
	tx       *Tx                    `json:"-"`
}

/**
* NewQuery
* @param model *Model, as string
* @return *Ql
**/
func NewQuery(model *Model, as string) *Ql {
	tp := SELECT
	if model.SourceField != "" {
		tp = DATA
	}
	maxRows := envar.GetInt("MAX_ROWS", 100)
	result := &Ql{
		Type:     tp,
		Froms:    make([]*From, 0),
		Joins:    make([]*Joins, 0),
		Selects:  make([]*Field, 0),
		Hiddens:  make([]string, 0),
		Wheres:   newWhere(),
		Details:  make(map[string]*Detail),
		Rollups:  make(map[string]*Detail),
		Calcs:    make(map[string]DataContext),
		GroupsBy: make([]*Field, 0),
		Havings:  newWhere(),
		OrdersBy: make([]*Orders, 0),
		Page:     0,
		Rows:     0,
		MaxRows:  maxRows,
		db:       model.db,
		tx:       nil,
	}
	result.addFrom(model, as)

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
* addFrom
* @param model *Model, as string
* @return *From
**/
func (s *Ql) addFrom(model *Model, as string) *From {
	result := model.from()
	result.As = as
	s.Froms = append(s.Froms, result)
	return result
}

/**
* findField
* @param field interface{}
* @return *Field
**/
func (s *Ql) findField(field interface{}) *Field {
	switch v := field.(type) {
	case string:
		return findFieldByStr(s.Froms, v)
	case *Agg:
		return findFieldByStr(s.Froms, v.Field)
	case *Field:
		return v
	default:
		return nil
	}
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
		s.Selects = make([]*Field, 0)
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
		f := s.findField(fld)
		if f == nil {
			continue
		}

		switch f.TypeColumn {
		case COLUMN:
			s.Selects = append(s.Selects, f)
		case AGG:
			s.Selects = append(s.Selects, f)
		case ATTRIB:
			s.Selects = append(s.Selects, f)
		case DETAIL:
			if f.From == nil {
				continue
			}

			details, err := s.db.GetModel(f.From.Key())
			if err != nil {
				continue
			}

			name := f.Name()
			detail, ok := details.Details[name]
			if !ok {
				continue
			}
			s.Details[name] = detail.setLimit(s.Page, s.Rows)
		case ROLLUP:
			if f.From == nil {
				continue
			}

			details, err := s.db.GetModel(f.From.Key())
			if err != nil {
				continue
			}

			name := f.Name()
			detail, ok := details.Rollups[name]
			if !ok {
				continue
			}
			s.Rollups[name] = detail.setLimit(s.Page, s.Rows)
		case CALC:
			if f.From == nil {
				continue
			}

			details, err := s.db.GetModel(f.From.Key())
			if err != nil {
				continue
			}

			name := f.Name()
			detail, ok := details.calcs[name]
			if !ok {
				continue
			}
			s.Calcs[name] = detail
		}
	}
	return s
}

/**
* Data
* @return *Ql
**/
func (s *Ql) Data(fields ...interface{}) *Ql {
	s.Select(fields...)
	s.Type = DATA
	return s
}

/**
* getDetails
* @param tx *Tx, data et.Json
**/
func (s *Ql) getDetails(tx *Tx, data et.Json) {
	for name, dtl := range s.Details {
		to := dtl.To
		model, err := s.db.GetModel(to.Key())
		if err != nil {
			return
		}

		ql := model.
			Select(dtl.Select...)
		for pk, fk := range dtl.Keys {
			val := data[pk]
			ql.Where(Eq(fk, val))
		}
		result, err := ql.LimitTx(tx, dtl.Page, dtl.Rows)
		if err != nil {
			return
		}

		data[name] = result.Result
	}
}

/**
* getRollups
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getRollups(tx *Tx, data et.Json) {
	for name, dtl := range s.Rollups {
		to := dtl.To
		model, err := s.db.GetModel(to.Name)
		if err != nil {
			return
		}

		flc := len(dtl.Select)
		ql := model.
			Select(dtl.Select...)
		for pk, fk := range dtl.Keys {
			val := data[pk]
			ql.Where(Eq(fk, val))
		}
		result, err := ql.LimitTx(tx, dtl.Page, dtl.Rows)
		if err != nil {
			return
		}

		if result.Count == 0 {
			if flc > 1 {
				data[name] = et.Json{}
			} else {
				data[name] = ""
			}
		} else if result.Count > 0 {
			if flc > 1 {
				data[name] = result.Result[0]
			} else {
				for _, v := range result.Result {
					data[name] = v
					break
				}
			}
		}
	}
}

/**
* getCalls
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getCalls(tx *Tx, data et.Json) {
	for _, call := range s.Calcs {
		call(tx, data)
	}
}

/**
* Join
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) join(tp TypeJoin, model *Model, as string, keys map[string]string) *Ql {
	from := s.addFrom(model, as)

	rKeys := make(map[string]string)
	for k, fk := range keys {
		field := s.findField(k)
		if field != nil {
			k = field.AS()
		}

		field = s.findField(fk)
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
	return s.join(JOIN, model, as, keys)
}

/**
* LeftJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) LeftJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(LEFT, model, as, keys)
}

/**
* RightJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) RightJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(RIGHT, model, as, keys)
}

/**
* FullJoin
* @param model *Model, as string, keys map[string]string
* @return *Ql
**/
func (s *Ql) FullJoin(model *Model, as string, keys map[string]string) *Ql {
	return s.join(FULL, model, as, keys)
}

/**
* Where
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) Where(condition *Condition) *Ql {
	fld := s.findField(condition.Field)
	if fld != nil {
		condition.Field = fld
	}
	s.Wheres.add(condition)
	return s
}

/**
* And
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) And(condition *Condition) *Ql {
	condition.Connector = AND
	s.Where(condition)
	return s
}

/**
* Or
* @param condition *Condition
* @return *Ql
**/
func (s *Ql) Or(condition *Condition) *Ql {
	condition.Connector = OR
	s.Where(condition)
	return s
}

/**
* GroupsBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) GroupBy(fields ...string) *Ql {
	for _, name := range fields {
		fld := s.findField(name)
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
		fld := s.findField(cnd.Field)
		if fld != nil {
			cnd.Field = fld
		}
		s.Havings.add(cnd)
	}
	return s
}

/**
* ordersBy
* @param asc bool, fields ...string
* @return *Ql
**/
func (s *Ql) ordersBy(asc bool, fields ...string) *Ql {
	for _, name := range fields {
		fld := s.findField(name)
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
* Hidden
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Hidden(fields ...string) *Ql {
	for _, name := range fields {
		fld := s.findField(name)
		if fld != nil {
			s.Hiddens = append(s.Hiddens, fld.AS())
		}
	}
	return s
}

/**
* ItExists
* @param data et.Json
* @return *Ql
**/
func (s *Ql) ItExists(data et.Json) *Ql {
	if len(s.Froms) == 0 {
		return s
	}

	model := s.Froms[0]
	if model == nil {
		return s
	}

	s.Type = EXISTS
	s.Wheres.ByPk(model, data)
	return s
}

/**
* AllTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Ql) AllTx(tx *Tx) (et.Items, error) {
	return s.db.Query(s)
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
* CountTx
* @param tx *Tx
* @return int, error
**/
func (s *Ql) CountTx(tx *Tx) (int, error) {
	result, err := s.OneTx(tx)
	if err != nil {
		return 0, err
	}

	count := result.Int("count")
	return count, nil
}

/**
* Count
* @return int, error
**/
func (s *Ql) Count() (int, error) {
	return s.CountTx(nil)
}

/**
* ExistsTx
* @param tx *Tx
* @return bool, error
**/
func (s *Ql) ExistsTx(tx *Tx) (bool, error) {
	s.Type = EXISTS
	result, err := s.OneTx(tx)
	if err != nil {
		return false, err
	}

	exists := result.Bool("exists")
	return exists, nil
}

/**
* Exists
* @return bool, error
**/
func (s *Ql) Exists() (bool, error) {
	return s.ExistsTx(nil)
}

/**
* Query
* @return Items, error
**/
func (s *Ql) Query(params et.Json) (et.Items, error) {
	return s.db.Query(s)
}
