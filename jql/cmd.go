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
		Wheres:        newWhere(),
		Data:          make([]et.Json, 0),
		beforeInserts: s.beforeInserts,
		beforeUpdates: s.beforeUpdates,
		beforeDeletes: s.beforeDeletes,
		afterInserts:  s.afterInserts,
		afterUpdates:  s.afterUpdates,
		afterDeletes:  s.afterDeletes,
		db:            s.db,
	}

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
* insert
* @return et.Items, error
**/
func (s *Cmd) insert() (et.Items, error) {
	result := et.Items{}
	for _, new := range s.Data {
		old := et.Json{}
		for _, fn := range s.beforeInserts {
			err := fn(s.tx, old, new)
			if err != nil {
				return et.Items{}, err
			}
		}

		if new.IsEmpty() {
			continue
		}

		result, err := s.db.Command(s)
		if err != nil {
			return et.Items{}, err
		}

		if !result.Ok {
			continue
		}

		new = result.First().Result
		for _, fn := range s.afterInserts {
			err := fn(s.tx, old, new)
			if err != nil {
				return et.Items{}, err
			}
		}

		result.Add(new)
	}

	return result, nil
}

/**
* update
* @return et.Items, error
**/
func (s *Cmd) update() (et.Items, error) {
	from := s.Model
	result := et.Items{}
	for _, data := range s.Data {
		current, err := from.
			Current(data).
			All()
		if err != nil {
			return et.Items{}, err
		}

		for _, old := range current.Result {
			new := old.Clone()
			for k, v := range data {
				new[k] = v
			}

			for _, fn := range s.beforeUpdates {
				err := fn(s.tx, old, new)
				if err != nil {
					return et.Items{}, err
				}
			}

			result, err := s.db.Command(s)
			if err != nil {
				return et.Items{}, err
			}

			if !result.Ok {
				continue
			}

			new = result.First().Result
			for _, fn := range s.afterUpdates {
				err := fn(s.tx, old, new)
				if err != nil {
					return et.Items{}, err
				}
			}

			result.Add(s.New)
		}
	}

	return result, nil
}

/**
* ExecTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Cmd) ExecTx(tx *Tx) (et.Items, error) {
	if s.db == nil {
		return et.Items{}, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	switch s.Type {
	case INSERT:
		return s.insert()
	case UPDATE:
		return s.update()
	case DELETE:
		return s.delete()
	case UPSERT:
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
* findField
* @param field interface{}
* @return *Field
**/
func (s *Cmd) findField(field interface{}) *Field {
	switch v := field.(type) {
	case string:
		return s.Model.findField(v)
	case *Field:
		return v
	default:
		return nil
	}
}

/**
* Where
* @param condition *Condition
* @return *Cmd
**/
func (s *Cmd) Where(condition *Condition) *Cmd {
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
* @return *Cmd
**/
func (s *Cmd) And(condition *Condition) *Cmd {
	condition.Connector = AND
	s.Where(condition)
	return s
}

/**
* Or
* @param condition *Condition
* @return *Cmd
**/
func (s *Cmd) Or(condition *Condition) *Cmd {
	condition.Connector = OR
	s.Where(condition)
	return s
}
