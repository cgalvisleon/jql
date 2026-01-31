package postgres

import (
	"fmt"

	"github.com/cgalvisleon/jql/jql"
)

func FieldAs(field *jql.Field) string {
	if field.From == nil {
		return field.As
	}

	return fmt.Sprintf(`%s.%s`, field.From.As, field.As)
}
