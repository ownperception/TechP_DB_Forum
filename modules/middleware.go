package modules

import (
	"database/sql"
	"log"
	"net/http"
)

func Check(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}

func UrlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type DbfuncHand func(w http.ResponseWriter, r *http.Request, db *sql.DB)
type DbHandler struct {
	Handle DbfuncHand
	Db     *sql.DB
}

func (h *DbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handle(w, r, h.Db)
}
