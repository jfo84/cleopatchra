package main

import (
	"net/http"

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

	http.HandleFunc("/", indexHandler)

	reposController := repos.NewController(db)
	http.HandleFunc("/repos", reposController.Get)

	pullsController := pulls.NewController(db)
	http.HandleFunc("/repos/{repoID}/pulls", pullsController.Get)

	pullController := pull.NewController(db)
	http.HandleFunc("/pulls/{pullID}", pullController.Get)

	repoController := repo.NewController(db)
	http.HandleFunc("/repos/{repoID}", repoController.Get)

	addr := ":7000"
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
