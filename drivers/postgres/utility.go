package postgres

import (
	"fmt"

	"github.com/cgalvisleon/jql/jdb"
)

func FieldAs(field *jdb.Field) string {
	if field.From == nil {
		return field.As
	}

	return fmt.Sprintf(`%s.%s`, field.From.As, field.As)
}
