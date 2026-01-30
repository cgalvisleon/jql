package jql

import "github.com/cgalvisleon/et/et"

func (s *Cmd) upsert() (et.Items, error) {
	if len(s.Data) == 0 {
		return et.Items{}, nil
	}

	data := s.Data[0]
	from := s.Model
	exists, err := from.
		ItExists(data).
		ItExists()
	if err != nil {
		return et.Items{}, err
	}

	if exists {
		s.Type = UPDATE
		s.Wheres.byPk(from, data)
		return s.update()
	}

	s.Type = INSERT
	return s.insert()
}
