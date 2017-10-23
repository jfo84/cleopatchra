package repos

import (
	"github.com/jfo84/cleopatchra/api/db"
	"net/http"
)

type ReposController struct {
	page int
	perPage int
	dbwrap *db.DBWrapper
}

func NewController(dbwrap *db.DBWrapper) *ReposController {
	return &ReposController{dbwrap:dbwrap}
}

func (rc *ReposController) Repos(w http.ResponseWriter, r *http.Request) {
	rc.dbwrap.GetRepos(rc.page, rc.perPage)
}