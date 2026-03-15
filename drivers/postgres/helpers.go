package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jql/jdb"
)

func FieldAs(field *jdb.Field) string {
	if field.From == nil {
		return field.As
	}

	result := field.From.As
	result = strs.Append(result, field.As, ".")
	return result
}
