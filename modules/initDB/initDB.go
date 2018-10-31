package initDB

import (
	mid "apiDB/modules"
	"database/sql"
	"io/ioutil"
)

func InitDB() *sql.DB {

	dbconf, err := ioutil.ReadFile("modules/initDB/dbconf")
	mid.Check(err)
	db, err := sql.Open("postgres", string(dbconf))
	mid.Check(err)

	freeDB, err := ioutil.ReadFile("modules/initDB/freeDB.sql")
	mid.Check(err)
	_, err = db.Exec(string(freeDB))
	mid.Check(err)

	initDB, err := ioutil.ReadFile("modules/initDB/initDB.sql")
	mid.Check(err)
	_, err = db.Exec(string(initDB))
	mid.Check(err)
	return db
}
