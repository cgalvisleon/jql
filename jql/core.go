package jql

/**
* initCore
* @param db *DB
* @return error
**/
func (s *DB) initCore() error {
	if err := defineCatalog(s); err != nil {
		return err
	}

	return nil
}
