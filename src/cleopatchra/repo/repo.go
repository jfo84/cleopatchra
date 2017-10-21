package repo

import (
	"db"
	"net/http"
)

type RepoController struct {
	id chan int
	db *DB
}

func NewController(id int, db *DB) *RepoController {
	return &RepoController{id:id, db:db}
}

func (rc *ReposController) Repo(w http.ResponseWriter, r *http.Request) {
	rc.db.GetRepo(u.id)
}

