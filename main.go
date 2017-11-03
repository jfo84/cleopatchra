package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/pull"
	"github.com/jfo84/cleopatchra/api/pulls"
	"github.com/jfo84/cleopatchra/api/repo"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func main() {
	db := db.OpenDb()
	r := mux.NewRouter().StrictSlash(true)
	s := r.PathPrefix("/repos").Subrouter()

	// reposController := repos.NewController(db)
	// s.HandleFunc("/", reposController.Get)

	repoController := repo.NewController(db)
	s.HandleFunc("/{repoID}/", repoController.Get)

	pullsController := pulls.NewController(db)
	s.HandleFunc("/{repoID}/pulls", pullsController.Get)

	pullController := pull.NewController(db)
	r.HandleFunc("/pulls/{pullID}", pullController.Get)

	// r.HandleFunc("/", indexHandler)

	addr := ":7000"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}
