package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/pull"
	"github.com/jfo84/cleopatchra/api/pulls"
	"github.com/jfo84/cleopatchra/api/repo"
	"github.com/jfo84/cleopatchra/api/repos"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func main() {
	db := db.OpenDb()
	r := mux.NewRouter()

	pullsController := pulls.NewController(db)
	r.HandleFunc("/repos/{repoID}/pulls", pullsController.Get)

	pullController := pull.NewController(db)
	r.HandleFunc("/pulls/{pullID}", pullController.Get)

	repoController := repo.NewController(db)
	r.HandleFunc("/repos/{repoID}", repoController.Get)

	reposController := repos.NewController(db)
	r.HandleFunc("/repos", reposController.Get)

	http.HandleFunc("/", indexHandler)

	addr := ":7000"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}
