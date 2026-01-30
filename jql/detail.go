package jql

type Detail struct {
	To              *From             `json:"to"`
	Keys            map[string]string `json:"keys"`
	Select          []interface{}     `json:"select"`
	OnDeleteCascade bool              `json:"on_delete_cascade"`
	OnUpdateCascade bool              `json:"on_update_cascade"`
}

/**
* newDetail
* @param to *Model, keys map[string]string, selecs []interface{}, onDeleteCascade, onUpdateCascade bool
* @return *Detail
**/
func newDetail(to *Model, keys map[string]string, selecs []interface{}, onDeleteCascade, onUpdateCascade bool) *Detail {
	return &Detail{
		To:              to.from(),
		Keys:            keys,
		Select:          selecs,
		OnDeleteCascade: onDeleteCascade,
		OnUpdateCascade: onUpdateCascade,
	}
}

type TypeJoin string

const (
	JOIN  TypeJoin = "join"
	LEFT  TypeJoin = "left"
	RIGHT TypeJoin = "right"
	FULL  TypeJoin = "full"
)

type Joins struct {
	Type TypeJoin
	To   *From
	Keys map[string]string
}

/**
* newJoins
* @param tp TypeJoin, from *From, keys map[string]string
* @return *Joins
**/
func newJoins(tp TypeJoin, from *From, keys map[string]string) *Joins {
	return &Joins{
		Type: tp,
		To:   from,
		Keys: keys,
	}
}
