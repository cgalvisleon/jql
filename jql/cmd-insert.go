package jql

import "github.com/cgalvisleon/et/et"

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
