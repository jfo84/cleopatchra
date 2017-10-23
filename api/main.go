package main

import (
	"net/http"
	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/repo"
	"github.com/jfo84/cleopatchra/api/repos"
	"github.com/jfo84/cleopatchra/api/pulls"
	"github.com/jfo84/cleopatchra/api/pull"
	"github.com/gorilla/mux"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func listenAndServe() {
	db := db.OpenDb()
	r := mux.NewRouter()

	reposController := repos.NewController(db)
	r.HandleFunc("/repos", reposController.Repos)

	repoController := repo.NewController(db)
	r.HandleFunc("/repo/{repoID}", repoController.Repo)

	pullsController := pulls.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls", pullsController.Pulls)

	pullController := pull.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls/{pullID}", pullController.Pull)

	http.HandleFunc("/", handleIndex)

	addr := ":7000"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}

func main() {
	listenAndServe()
}
