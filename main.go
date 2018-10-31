package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	hand "github.com/ownperception/TechP_DB_Forum/modules/handlers"
	idb "github.com/ownperception/TechP_DB_Forum/modules/initDB"
	mid "github.com/ownperception/TechP_DB_Forum/modules/middlefunc"
)

var db *sql.DB

func init() {
	db = idb.InitDB()
}

func main() {
	defer db.Close()

	r := mux.NewRouter()
	r.Handle("/forum/create", &mid.DbHandler{Handle: hand.ForumCreate, Db: db}).Methods("POST")
	r.Handle("/forum/{id}/create", &mid.DbHandler{Handle: hand.ForumTrCreate, Db: db}).Methods("POST")
	r.Handle("/forum/{id}/{method:details|threads|users}", &mid.DbHandler{Handle: hand.ForumStat, Db: db}).Methods("GET")
	r.Handle("/post/{id}/details", &mid.DbHandler{Handle: hand.PostInfo, Db: db}).Methods("POST", "GET")
	r.Handle("/service/{method:clear|status}", &mid.DbHandler{Handle: hand.Service, Db: db}).Methods("POST", "GET")
	r.Handle("/thread/{slug_or_id}/{method:details|posts}", &mid.DbHandler{Handle: hand.ThreadGet, Db: db}).Methods("GET")
	r.Handle("/thread/{slug_or_id}/{method:create|details|vote}", &mid.DbHandler{Handle: hand.ThreadPost, Db: db}).Methods("POST")
	r.Handle("/user/{nickname}/{method:create|profile}", &mid.DbHandler{Handle: hand.UserPost, Db: db}).Methods("POST")
	r.Handle("/user/{nickname}/profile", &mid.DbHandler{Handle: hand.UserGet, Db: db}).Methods("GET")

	router := http.NewServeMux()
	router.Handle("/", mid.UrlMiddleware(r))

	fmt.Println("starting server at :5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
