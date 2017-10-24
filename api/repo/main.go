package repo

import (
	"net/http"

	"github.com/jfo84/cleopatchra/api/db"
)

// Controller - For re-use of *db.Wrapper
type Controller struct {
	dbwrap *db.Wrapper
}

// NewController is a constructor for initializing with a *db.Wrapper
func NewController(dbwrap *db.Wrapper) *Controller {
	return &Controller{dbwrap: dbwrap}
}

// Get writes the controller's model values with the http.ResponseWriter
func (rc *Controller) Get(w http.ResponseWriter, r *http.Request) {
	rc.dbwrap.GetRepo(w, r)
}
