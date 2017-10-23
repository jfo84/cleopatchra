package repo

import (
	"github.com/jfo84/cleopatchra/api/db"
	"net/http"
)

type RepoController struct {
	id int
	dbwrap *db.DBWrapper
}

func NewController(dbwrap *db.DBWrapper) *RepoController {
	return &RepoController{dbwrap:dbwrap}
}

func (rc *RepoController) Repo(w http.ResponseWriter, r *http.Request) {
	rc.dbwrap.GetRepo(rc.id)
}