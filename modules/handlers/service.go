package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	mod "github.com/ownperception/TechP_DB_Forum/modules/middlefunc"
)

func Service(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)

	switch vars["method"] {
	case "clear":
		if r.Method != http.MethodPost {
			fmt.Println("StatusMethodNotAllowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		_, err := db.Exec("TRUNCATE TABLE author,post,thread,forum,vote;")
		mod.Check(err)
	case "status":
		if r.Method != http.MethodGet {
			fmt.Println("StatusMethodNotAllowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		row := db.QueryRow("select cf,cp,ct,ca " +
			"from (select count(*) cf from forum) s1," +
			"(select count(*) ca from author) s4," +
			"(select count(*) ct from thread) s3," +
			"(select count(*) cp from post) s2;")
		var f, p, t, u int
		err := row.Scan(&f, &p, &t, &u)
		mod.Check(err)
		arr := map[string]int{
			"forum":  f,
			"post":   p,
			"thread": t,
			"user":   u,
		}
		data, err := json.Marshal(arr)
		mod.Check(err)
		w.Write(data)
	}
	w.WriteHeader(http.StatusOK)
}
