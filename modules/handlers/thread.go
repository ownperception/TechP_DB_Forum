package handlers

import (
	types "apiDB/models"
	mod "apiDB/modules"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ThreadGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]
	id_flag := mod.FlagSlugOrId(slug_or_id)

	var data []byte
	tr := types.Thread{}
	reqstring := fmt.Sprintf("select * from thread where %s ;", id_flag)

	err := db.QueryRow(reqstring).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
	if err == sql.ErrNoRows {
		msg := "Can't find thread by slug_or_id: " + slug_or_id
		data, _ = json.Marshal(types.Error{Message: msg})
		w.WriteHeader(http.StatusNotFound)
	} else {
		switch vars["method"] {
		case "details":
			data, _ = json.Marshal(tr)
			w.WriteHeader(http.StatusOK)
		case "posts":
			params := map[string]string{
				"limit": "",
				"since": "",
				"desc":  "false",
				"sort":  "flat",
			}
			mod.ParsUrl(r, &params)

			flags := map[string]string{
				"sortflag":  "",
				"limitflag": "",
				"sinceflag": "",
			}

			switch params["sort"] {
			case "flat":
				mod.ReqFlatFlags(params, &flags)

				reqstring := fmt.Sprintf("select * from post where thread = $1 %s order by id %s %s;", flags["sinceflag"], flags["sortflag"], flags["limitflag"])
				log.Println(reqstring)

				rows, err := db.Query(reqstring, tr.Id)
				defer rows.Close()
				mod.Check(err)
				ps := []types.Post{}

				for rows.Next() {
					p := types.Post{}
					err := rows.Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
					mod.Check(err)
					ps = append(ps, p)
				}
				data, _ = json.Marshal(ps)
			case "tree":
				mod.ReqTreeFlags(params, &flags)

				reqstring := fmt.Sprintf("with recursive post_tree(id,path) as( "+
					"select p.id,array_append('{}'::bigint[], id) as arr_id "+
					"from post p "+
					"where p.parent = 0 and p.thread = $1 "+

					"union all "+

					"select p.id, array_append(path, p.id) from post p "+
					"join post_tree pt on p.parent = pt.id "+
					") "+
					"select p.id,p.author,p.created,p.message,p.forum,p.thread,p.isedited,p.parent from post_tree pt join post p on p.id = pt.id %s %s %s;", flags["sinceflag"], flags["sortflag"], flags["limitflag"])
				rows, err := db.Query(reqstring, tr.Id)
				defer rows.Close()
				mod.Check(err)
				ps := []types.Post{}

				for rows.Next() {
					p := types.Post{}
					err := rows.Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
					mod.Check(err)
					ps = append(ps, p)
				}
				data, _ = json.Marshal(ps)

			case "parent_tree":
				mod.ReqParTreeFlags(params, &flags)

				reqstring := fmt.Sprintf("select p.id,p.author,p.created,p.message,p.forum,p.thread,p.isedited,p.parent from (with recursive post_tree(id,path) as( "+
					"select p.id,array_append('{}'::bigint[], p.id) as arr_id "+
					"from post p "+
					"where p.parent = 0 and p.thread = $1 "+

					"union all "+

					"select p.id, array_append(path, p.id) from post p "+
					"join post_tree pt on p.parent = pt.id "+
					") "+
					"select post_tree.id as id,path, dense_rank() over (order by path[1] %s ) as r from post_tree %s ) as pt join post p on p.id = pt.id %s %s;", flags["descflag"], flags["sinceflag"], flags["limitflag"], flags["sortflag"])
				log.Println(reqstring)
				rows, err := db.Query(reqstring, tr.Id)
				defer rows.Close()
				mod.Check(err)
				ps := []types.Post{}

				for rows.Next() {
					p := types.Post{}
					err := rows.Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
					mod.Check(err)
					ps = append(ps, p)
				}
				data, _ = json.Marshal(ps)
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	w.Write(data)
}

func ThreadPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	slug_or_id := vars["slug_or_id"]
	id_flag := mod.FlagSlugOrId(slug_or_id)

	var data []byte

	tr := types.Thread{}
	reqstring := fmt.Sprintf("select * from thread where %s ;", id_flag)

	err := db.QueryRow(reqstring).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
	if err == sql.ErrNoRows {
		msg := "Can't find thread by slug_or_id: " + slug_or_id
		data, _ = json.Marshal(types.Error{Message: msg})
		w.WriteHeader(http.StatusNotFound)
	} else {
		mod.Check(err)
		switch vars["method"] {
		case "details":
			params, err := mod.Jsonparams(r)

			args := []string{}
			for key, val := range params {
				if val != "" {
					args = append(args, key+" = '"+val+"'")
				}
			}

			var reqstring string
			if len(args) != 0 {
				reqstring = "update thread set " + strings.Join(args, ",") + fmt.Sprintf(" where id = %d returning *;", tr.Id)
				log.Println(reqstring)

				err = db.QueryRow(reqstring).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
				mod.Check(err)
			}

			data, _ = json.Marshal(tr)
			w.WriteHeader(http.StatusOK)

		case "create":
			ps := []types.Post{}
			reqstart := "insert into post (thread,forum,author,message,parent) values "
			reqEnd := " returning *;"
			Parr := []string{}

			arrayparams, _ := mod.GetJsonArrayPost(r)
			for _, post := range arrayparams {
				poststring := ""
				if post.Parent != 0 {
					var parThread string
					reqstring = fmt.Sprintf("select thread from post where id = '%d'", post.Parent)

					err := db.QueryRow(reqstring).Scan(&parThread)
					if err != nil {
						if err == sql.ErrNoRows {
							msg := "Can't find parent post "
							data, _ = json.Marshal(types.Error{Message: msg})
							w.WriteHeader(http.StatusConflict)
							w.Write(data)
							return
						}
					}

					if parThread != strconv.Itoa(tr.Id) {
						msg := "Parent post was created in another thread"
						data, _ = json.Marshal(types.Error{Message: msg})
						w.WriteHeader(http.StatusConflict)
						w.Write(data)
						return
					}
					poststring = fmt.Sprintf("(%d,'%s','%s','%s',%d)", tr.Id, tr.Forum, post.Author, post.Message, post.Parent)

				} else {
					poststring = fmt.Sprintf("(%d,'%s','%s','%s', 0 )", tr.Id, tr.Forum, post.Author, post.Message)
				}
				Parr = append(Parr, poststring)
			}
			reqstring := reqstart + strings.Join(Parr, ",") + reqEnd

			if len(Parr) > 0 {
				rows, err := db.Query(reqstring)

				if err != nil {
					switch err.Error() {
					case "pq: insert or update on table \"post\" violates foreign key constraint \"post_author_fkey\"":
						msg := "No author"
						data, _ = json.Marshal(types.Error{Message: msg})
						w.WriteHeader(http.StatusNotFound)
						w.Write(data)
						return
					default:
						log.Println(reqstring)
						log.Fatalln(err)
					}
				} else {
					defer rows.Close()
					for rows.Next() {
						p := types.Post{}
						rows.Scan(&p.Id, &p.Author, &p.Created, &p.Message, &p.Forum, &p.Thread, &p.IsEdited, &p.Parent)
						ps = append(ps, p)
					}
				}
			}
			data, _ = json.Marshal(ps)
			w.WriteHeader(http.StatusCreated)

		case "vote":
			vote, err := mod.GetJsonVote(r)

			res, err := db.Exec("update vote set voice = $2 where nickname = $1 and thread = $3;", vote.Nickname, vote.Voice, tr.Id)
			if err != nil {
				if err.Error() == "pq: insert or update on table \"vote\" violates foreign key constraint \"vote_nickname_fkey\"" {
					msg := "Can't find user : " + vote.Nickname
					data, _ = json.Marshal(types.Error{Message: msg})
					w.WriteHeader(http.StatusNotFound)
				} else {
					log.Fatalln(err)
				}
			}

			num, _ := res.RowsAffected()
			if num == 0 {
				_, err := db.Exec("insert into vote (nickname,voice,thread) values ($1,$2,$3);", vote.Nickname, vote.Voice, tr.Id)
				if err != nil {
					if err.Error() == "pq: insert or update on table \"vote\" violates foreign key constraint \"vote_nickname_fkey\"" {
						msg := "Can't find user : " + vote.Nickname
						data, _ = json.Marshal(types.Error{Message: msg})
						w.WriteHeader(http.StatusNotFound)
					} else {
						log.Fatalln(err)
					}
				}
			}
			reqstring = fmt.Sprintf("select * from thread where id = %d;", tr.Id)
			err = db.QueryRow(reqstring).Scan(&tr.Id, &tr.Author, &tr.Created, &tr.Forum, &tr.Title, &tr.Message, &tr.Slug, &tr.Votes)
			mod.Check(err)
			data, _ = json.Marshal(tr)
			w.WriteHeader(http.StatusOK)
		}
	}
	w.Write(data)
}
