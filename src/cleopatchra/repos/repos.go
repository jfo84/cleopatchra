package repos

import (
	"db"
	"net/http"
)

type ReposController struct {
	page chan int
	perPage chan int
	db *DB
}

func NewController(db *DB) *ReposController {
	return &ReposController{db:db}
}

func (u *ReposController) GetRepos(w http.ResponseWriter, r *http.Request) {
	u.db.GetRepos(u.page, u.perPage)
}