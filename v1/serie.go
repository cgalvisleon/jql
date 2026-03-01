package jql

import "github.com/cgalvisleon/jql/jdb"

/**
* SetSerie
* @param tag, format string, value int
* @return error
**/
func SetSerie(tag, format string, value int) error {
	return jdb.SetSerie(tag, format, value)
}

/**
* GetSerie
* @param tag string
* @return (int, string, error)
**/
func GetSerie(tag string) (int, string, error) {
	return jdb.GetSerie(tag)
}

/**
* DeleteSerie
* @param tag string
* @return error
**/
func DeleteSerie(tag string) error {
	return jdb.DeleteSerie(tag)
}
