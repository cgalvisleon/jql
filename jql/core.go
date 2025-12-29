package jql

func (s *DB) initCore() error {
	if err := defineModel(s); err != nil {
		return err
	}

	return nil
}
