package pulls

import (
	"db"
	"net/http"
)

type PullsController struct {
	repoID chan string
	page chan int
	perPage chan int
	db *DB
}

func NewController(db *DB) *PullsController {
	return &PullsController{db:db}
}

func (pc *PullsController) GetPulls(w http.ResponseWriter, r *http.Request) {
	pc.db.GetPulls(pc.repoID, pc.page, pc.perPage)
}