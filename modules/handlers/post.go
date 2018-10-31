package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	types "github.com/ownperception/TechP_DB_Forum/apiDB/models"
	mod "github.com/ownperception/TechP_DB_Forum/apiDB/modules/middlefunc"
)

func PostInfo(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	id := vars["id"]

	var data []byte

	row := db.QueryRow("select * from post where id = $1;", id)
	p := types.Post{}
	err := row.Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "Can't find post by id: " + id
			data, _ = json.Marshal(types.Error{Message: msg})
			w.WriteHeader(http.StatusNotFound)
		}
	} else {

		if r.Method == http.MethodPost {
			params, err := mod.Jsonparams(r)

			args := []string{}
			for key, val := range params {
				if val != "" {
					args = append(args, key+" = '"+val+"'")
				}
			}

			var reqstring string
			if len(args) != 0 && params["message"] != p.Message {
				reqstring = "update post set " + strings.Join(args, ",") + fmt.Sprintf(", isedited = true where id = %s returning *;", id)
				err = db.QueryRow(reqstring).Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
				mod.Check(err)
			}

			data, _ = json.Marshal(p)

		} else {
			params := map[string]string{
				"related": "",
			}
			mod.ParsUrl(r, &params)

			arr := map[string]interface{}{
				"post": p,
			}
			log.Println(params["related"])
			if strings.Contains(params["related"], "user") {
				a := types.Author{}
				err := db.QueryRow("select * from author where nickname = $1;", p.Author).Scan(&a.Id, &a.Fullname, &a.Nickname, &a.Email, &a.About)
				mod.Check(err)
				arr["author"] = a
			}
			if strings.Contains(params["related"], "thread") {
				tr := types.Thread{}
				err = db.QueryRow("select * from thread where id = $1;", p.Thread).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
				mod.Check(err)
				arr["thread"] = tr
			}
			if strings.Contains(params["related"], "forum") {
				f := types.Forum{}
				err := db.QueryRow("select * from forum where slug = $1;", p.Forum).Scan(&f.Slug, &f.Title, &f.Author, &f.Posts, &f.Threads)
				mod.Check(err)
				arr["forum"] = f
			}

			data, _ = json.Marshal(arr)
			w.WriteHeader(http.StatusOK)

		}
	}
	w.Write(data)
}
