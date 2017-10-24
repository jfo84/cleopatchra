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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func main() {
	db := db.OpenDb()
	r := mux.NewRouter()

	reposController := repos.NewController(db)
	r.HandleFunc("/repos", reposController.Get)

	repoController := repo.NewController(db)
	r.HandleFunc("/repo/{repoID}", repoController.Get)

	pullsController := pulls.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls", pullsController.Get)

	pullController := pull.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls/{pullID}", pullController.Get)

	http.HandleFunc("/", handleIndex)

	addr := ":7000"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}
