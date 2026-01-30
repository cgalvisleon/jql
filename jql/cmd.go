package jql

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
)

type TypeCommand string

const (
	INSERT TypeCommand = "insert"
	UPDATE TypeCommand = "update"
	DELETE TypeCommand = "delete"
	UPSERT TypeCommand = "upsert"
)

type Cmd struct {
	Type          TypeCommand       `json:"type"`
	Model         *Model            `json:"model"`
	Wheres        *Wheres           `json:"wheres"`
	Data          []et.Json         `json:"data"`
	New           et.Json           `json:"new"`
	IsDebug       bool              `json:"is_debug"`
	beforeInserts []TriggerFunction `json:"-"`
	beforeUpdates []TriggerFunction `json:"-"`
	beforeDeletes []TriggerFunction `json:"-"`
	afterInserts  []TriggerFunction `json:"-"`
	afterUpdates  []TriggerFunction `json:"-"`
	afterDeletes  []TriggerFunction `json:"-"`
	tx            *Tx               `json:"-"`
	db            *DB               `json:"-"`
}

/**
* Serialize
* @return []byte, error
**/
func (s *Cmd) serialize() ([]byte, error) {
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
func (s *Cmd) ToJson() et.Json {
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
* newCommand
* @param model *Model, cmd TypeCommand
* @return *Cmd
**/
func newCommand(s *Model, cmd TypeCommand) *Cmd {
	result := &Cmd{
		Type:          cmd,
		Model:         s,
		Data:          make([]et.Json, 0),
		New:           et.Json{},
		beforeInserts: s.beforeInserts,
		beforeUpdates: s.beforeUpdates,
		beforeDeletes: s.beforeDeletes,
		afterInserts:  s.afterInserts,
		afterUpdates:  s.afterUpdates,
		afterDeletes:  s.afterDeletes,
		db:            s.db,
	}
	result.Wheres = newWhere(result)

	return result
}

/**
* setTx
* @param tx *Tx
* @return *Cmd
**/
func (s *Cmd) setTx(tx *Tx) *Cmd {
	s.tx = tx
	return s
}

/**
* ExecTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Cmd) ExecTx(tx *Tx) (et.Items, error) {
	if s.DB == nil {
		return et.Items{}, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	switch s.Type {
	case TypeInsert:
		return s.insert()
	case TypeUpdate:
		return s.update()
	case TypeDelete:
		return s.delete()
	case TypeUpsert:
		return s.upsert()
	default:
		return et.Items{}, fmt.Errorf("invalid command: %s", s.Type)
	}
}

/**
* Exec
* @return et.Items, error
**/
func (s *Cmd) Exec() (et.Items, error) {
	return s.ExecTx(nil)
}

/**
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Cmd) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.ExecTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* One
* @return et.Item, error
**/
func (s *Cmd) One() (et.Item, error) {
	return s.OneTx(nil)
}

/**
* Where
* @param condition *Condition
* @return *Cmd
**/
func (s *Cmd) Where(condition *Condition) *Cmd {
	s.Wheres.Add(condition)
	return s
}

/**
* WhereByWhere
* @param where *Wheres
* @return *Cmd
**/
func (s *Cmd) WhereByJson(where []et.Json) *Cmd {
	for _, w := range where {
		field := w.String("field")
		value := w.Get("value")
		operator := Operator(w.String("operator"))
		condition := condition(field, value, operator)
		s.Wheres.Add(condition)
	}
	return s
}

/**
* WhereByPrimaryKeys
* @param data et.Json
* @return *Cmd
**/
func (s *Cmd) WhereByPrimaryKeys(data et.Json) *Cmd {
	s.Wheres = newWhere(s)
	for _, col := range s.Model.PrimaryKeys {
		val := data[col]
		if val == nil {
			continue
		}
		s.Where(Eq(col, val))
	}

	return s
}

/**
* And
* @param condition *Condition
* @return *Cmd
**/
func (s *Cmd) And(condition *Condition) *Cmd {
	s.Wheres.Add(condition)
	return s
}

/**
* Or
* @param condition *Condition
* @return *Cmd
**/
func (s *Cmd) Or(condition *Condition) *Cmd {
	s.Wheres.Add(condition)
	return s
}
