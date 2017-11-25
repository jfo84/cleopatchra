package unmarshalling

import (
	"github.com/jfo84/cleopatchra/api/exports"
)

/*
This package exists for edge cases where the struct we want to build for exports does not
fit cleanly into the json.Unmarshal interface.

With Pull, for example, we need to side-load comment internal ID's. However, the GitHub JSON
contains a "comments" key with an int value (# of comments)
*/

// Pull represents the unmarshalling interface for a GitHub pull request
type Pull struct {
	ID          int           `json:"id"`
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	Merged      bool          `json:"merged"`
	NumComments int           `json:"comments"`
	User        *exports.User `json:"user"`
	Repo        *exports.Repo `json:"repo"`
}
