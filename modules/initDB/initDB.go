package initDB

import (
	"database/sql"
	"io/ioutil"

	mid "github.com/ownperception/TechP_DB_Forum/modules/middlefunc"
)

func InitDB() *sql.DB {

	dbconf, err := ioutil.ReadFile("github.com/ownperception/TechP_DB_Forum/modules/initDB/dbconf")
	mid.Check(err)
	db, err := sql.Open("postgres", string(dbconf))
	mid.Check(err)

	freeDB, err := ioutil.ReadFile("github.com/ownperception/TechP_DB_Forum/modules/initDB/freeDB.sql")
	mid.Check(err)
	_, err = db.Exec(string(freeDB))
	mid.Check(err)

	initDB, err := ioutil.ReadFile("github.com/ownperception/TechP_DB_Forum/modules/initDB/initDB.sql")
	mid.Check(err)
	_, err = db.Exec(string(initDB))
	mid.Check(err)
	return db
}
