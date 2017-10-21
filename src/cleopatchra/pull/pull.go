package pull

import (
	"db"
	"net/http"
)

type PullController struct {
	id chan int
	db *DB
}

func NewController(db *DB) *PullController {
	return &PullController{db:db}
}

func (pc *PullController) GetPull(w http.ResponseWriter, r *http.Request) {
	pc.db.GetPull(pc.id)
}