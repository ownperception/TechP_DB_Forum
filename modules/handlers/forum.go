package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	types "github.com/ownperception/TechP_DB_Forum/apiDB/models"
	mod "github.com/ownperception/TechP_DB_Forum/apiDB/modules"
)

func ForumCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	params, err := mod.Jsonparams(r)

	slug := params["slug"]
	title := params["title"]
	user := params["user"]
	var data []byte

	err = db.QueryRow("select nickname from author where \"nickname\" = $1;", user).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "Can't find user by nickname: " + user
			data, _ = json.Marshal(types.Error{Message: msg})
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Fatal(err)
		}
	} else {
		f := types.Forum{}
		err = db.QueryRow("insert into forum values ($1,$2,$3) returning *;", slug, title, user).Scan(&f.Slug, &f.Title, &f.Author, &f.Posts, &f.Threads)
		if err != nil {
			switch err.Error() {
			case "pq: duplicate key value violates unique constraint \"forum_pkey\"":
				f := types.Forum{}
				err := db.QueryRow("select * from forum where \"slug\"= $1;", slug).Scan(&f.Slug, &f.Title, &f.Author, &f.Posts, &f.Threads)
				mod.Check(err)
				data, _ = json.Marshal(f)
				w.WriteHeader(http.StatusConflict)
			default:
				log.Fatal(err)
			}
		} else {
			data, _ = json.Marshal(f)
			w.WriteHeader(http.StatusCreated)
		}
	}
	w.Write(data)
}

func ForumTrCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	params, err := mod.Jsonparams(r)

	forum := vars["id"]
	slug := params["slug"]
	title := params["title"]
	message := params["message"]
	created := params["created"]
	author := params["author"]
	var data []byte

	err = db.QueryRow("select nickname from author where \"nickname\" = $1;", author).Scan(&author)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "Can't find author by nickname: " + author
			data, _ = json.Marshal(types.Error{Message: msg})
			w.WriteHeader(http.StatusNotFound)
			w.Write(data)
			return
		} else {
			log.Fatal(err)
		}
	}
	tr := types.Thread{}

	if slug != "" {
		err = db.QueryRow("select * from thread where \"slug\" = $1;", slug).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
		if err != sql.ErrNoRows {
			data, _ = json.Marshal(tr)
			w.WriteHeader(http.StatusConflict)
			w.Write(data)
			return
		}
	}

	err = db.QueryRow("select slug from forum where \"slug\"= $1;", forum).Scan(&forum)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "Can't find forum by slug: " + forum
			data, _ = json.Marshal(types.Error{Message: msg})
			w.WriteHeader(http.StatusNotFound)
			w.Write(data)
			return
		} else {
			log.Fatal(err)
		}
	}

	if created != "" {
		err = db.QueryRow("insert into thread (author,forum,title,message,slug,created) values ($1,$2,$3,$4,$5,$6) returning *;", author, forum, title, message, slug, created).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
	} else {
		err = db.QueryRow("insert into thread (author,forum,title,message,slug) values ($1,$2,$3,$4,$5) returning *;", author, forum, title, message, slug).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
	}
	if err != nil {
		log.Fatal(err)
	} else {
		data, _ = json.Marshal(tr)
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}

}

func ForumStat(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	slug := vars["id"]
	var data []byte
	f := types.Forum{}

	err := db.QueryRow("select * from forum where slug = $1 ;", slug).Scan(&f.Slug, &f.Title, &f.Author, &f.Posts, &f.Threads)
	if err == sql.ErrNoRows {
		msg := "Can't find forum by slug: " + slug
		data, _ = json.Marshal(types.Error{Message: msg})
		w.WriteHeader(http.StatusNotFound)
	} else {
		if vars["method"] == "details" {
			data, _ = json.Marshal(f)
			w.WriteHeader(http.StatusOK)
		} else {
			params := map[string]string{
				"limit": "",
				"since": "",
				"desc":  "false",
			}
			mod.ParsUrl(r, &params)

			flags := map[string]string{
				"sortflag":  "",
				"limitflag": "",
				"sinceflag": "",
			}

			switch vars["method"] {
			case "users":
				mod.ReqUsersFlags(params, &flags)

				reqstring := fmt.Sprintf("select distinct id, fullname, nickname, email, about from ("+
					"select distinct author "+
					"from post "+
					"where forum = $1 %s"+
					"union all "+
					"select distinct author "+
					"from thread "+
					"where forum = $1 %s "+
					") as sel join author a on a.nickname = sel.author "+
					"order by nickname %s %s;", flags["sinceflag"], flags["sinceflag"], flags["sortflag"], flags["limitflag"])
				log.Println(reqstring)
				rows, err := db.Query(reqstring, slug)
				defer rows.Close()
				mod.Check(err)

				usrs := []types.Author{}
				for rows.Next() {
					usr := types.Author{}
					err := rows.Scan(&usr.Id, &usr.Fullname, &usr.Nickname, &usr.Email, &usr.About)
					mod.Check(err)
					usrs = append(usrs, usr)
				}
				data, _ = json.Marshal(usrs)
				w.WriteHeader(http.StatusOK)
			case "threads":
				mod.ReqThreadsFlags(params, &flags)

				reqstring := fmt.Sprintf("select * from (select * from thread where forum = $1 order by created %s) s1 %s %s;", flags["sortflag"], flags["sinceflag"], flags["limitflag"])
				log.Println(reqstring)
				rows, err := db.Query(reqstring, slug)
				defer rows.Close()
				mod.Check(err)

				trs := []types.Thread{}
				for rows.Next() {
					tr := types.Thread{}
					err := rows.Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
					mod.Check(err)
					trs = append(trs, tr)
				}
				data, _ = json.Marshal(trs)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	w.Write(data)
}
