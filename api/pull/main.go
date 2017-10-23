package pull

import (
	"github.com/jfo84/cleopatchra/api/db"
	"net/http"
)

type PullController struct {
	id int
	dbwrap *db.DBWrapper
}

func NewController(dbwrap *db.DBWrapper) *PullController {
	return &PullController{dbwrap:dbwrap}
}

func (pc *PullController) Pull(w http.ResponseWriter, r *http.Request) {
	pc.dbwrap.GetPull(pc.id)
}