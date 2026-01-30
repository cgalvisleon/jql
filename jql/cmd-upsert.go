package jql

import "github.com/cgalvisleon/et/et"

func (s *Cmd) upsert() (et.Items, error) {
	if s.Model == nil {
		return et.Items{}, nil
	}

	from := s.Model
	for _, data := range s.Data {
		exists, err := from.
			ItExists(data).
			ItExists()
		if err != nil {
			return et.Items{}, err
		}

		if exists {
			s.Type = UPDATE
			s.WhereByPrimaryKeys(data)
			return s.update()
		}
	}

	if exists {
		s.Type = TypeUpdate
		s.WhereByPrimaryKeys(data)
		return s.update()
	}

	s.Type = TypeInsert
	return s.insert()
}
