package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	types "github.com/ownperception/TechP_DB_Forum/models"
	mod "github.com/ownperception/TechP_DB_Forum/modules/middlefunc"
)

func UserPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	params, err := mod.Jsonparams(r)

	nickname := vars["nickname"]
	fullname := params["fullname"]
	email := params["email"]
	about := params["about"]

	var data []byte

	switch vars["method"] {
	case "create":
		usr := types.Author{}

		row := db.QueryRow("insert into author (nickname,email,fullname,about) values ($1,$2,$3,$4) returning fullname,nickname,email,about;", nickname, email, fullname, about)
		err = row.Scan(&usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"author_nickname_key\"" || err.Error() == "pq: duplicate key value violates unique constraint \"author_email_key\"" {
				rows, err := db.Query("select * from author where \"nickname\"= $1 or email = $2;", nickname, email)
				mod.Check(err)

				arr := []types.Author{}
				for rows.Next() {
					err = rows.Scan(&usr.Id, &usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)
					mod.Check(err)
					arr = append(arr, usr)
				}
				data, err = json.Marshal(arr)
				w.WriteHeader(http.StatusConflict)

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			data, err = json.Marshal(usr)
			mod.Check(err)
			w.WriteHeader(http.StatusCreated)
		}

	case "profile":
		usr := types.Author{}

		reqparams := make(map[string]string)
		reqstart := "update author set "
		reqend := "where \"nickname\" = $1 returning *;"
		if email != "" {
			reqparams["email"] = email
		}
		if fullname != "" {
			reqparams["fullname"] = fullname
		}
		if about != "" {
			reqparams["about"] = about
		}
		if len(reqparams) == 0 {
			row := db.QueryRow("select * from author where \"nickname\" = $1;", nickname)
			err := row.Scan(&usr.Id, &usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)
			mod.Check(err)
			data, _ = json.Marshal(usr)
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
		idx := 0
		for key, val := range reqparams {
			if idx == 0 {
				reqstart = fmt.Sprintf(reqstart+"%s = '%s'", key, val)
			} else {
				reqstart = fmt.Sprintf(reqstart+", %s = '%s'", key, val)
			}
			idx++
		}
		reqstring := reqstart + reqend

		row := db.QueryRow(reqstring, nickname)
		err := row.Scan(&usr.Id, &usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)

		if err != nil {
			var msg string
			switch err.Error() {
			case "pq: duplicate key value violates unique constraint \"author_nickname_key\"":
				msg = "This nickname is already registered"
			case "pq: duplicate key value violates unique constraint \"author_email_key\"":
				row := db.QueryRow("select nickname from author where email = $1;", email)
				err = row.Scan(&usr.Nickname)
				mod.Check(err)
				msg = "This email is already registered by user: " + usr.Nickname
			default:
				msg := "Can't find user by nickname: " + nickname
				data, _ := json.Marshal(types.Error{Message: msg})
				w.WriteHeader(http.StatusNotFound)
				w.Write(data)
				return
			}

			errMsg := types.Error{Message: msg}
			data, _ = json.Marshal(errMsg)
			w.WriteHeader(http.StatusConflict)

		} else {
			data, _ = json.Marshal(usr)
			w.WriteHeader(http.StatusOK)
		}
	}
	w.Write(data)
}

func UserGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	usr := types.Author{}
	row := db.QueryRow("select * from author where \"nickname\" = $1;", nickname)
	err := row.Scan(&usr.Id, &usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)
	if err != nil {
		msg := "Can't find user by nickname: " + usr.Nickname
		data, _ := json.Marshal(types.Error{Message: msg})
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)
		return
	}
	data, _ := json.Marshal(usr)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
