package pulls

import (
	"db"
	"net/http"
)

type PullsController struct {
	repoID *string
	page int
	perPage int
	dbwrap *db.DBWrapper
}

func NewController(dbwrap *db.DBWrapper) *PullsController {
	return &PullsController{dbwrap:dbwrap}
}

func (pc *PullsController) Pulls(w http.ResponseWriter, r *http.Request) {
	pc.dbwrap.GetPulls(pc.repoID, pc.page, pc.perPage)
}