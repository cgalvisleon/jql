package jql

var models *Model

/**
* defineModel
* @param db *DB
* @return error
**/
func defineModel(db *DB) error {
	if models != nil {
		return nil
	}

	var err error
	models, err = db.NewModel("core", "models", 1)
	if err != nil {
		return err
	}
	models.defineCreatedAtField()
	models.defineUpdatedAtField()
	models.DefineColumn("name", TEXT, "")
	models.DefineColumn("version", INT, 0)
	models.DefineColumn("definition", BYTES, []byte{})
	models.DefinePrimaryKeys("name")
	models.IsCore = true
	if err = models.Init(); err != nil {
		return err
	}

	return nil
}

/**
* initCore
* @param db *DB
* @return error
**/
func (s *DB) initCore() error {
	if err := defineModel(s); err != nil {
		return err
	}

	return nil
}
