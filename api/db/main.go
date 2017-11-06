package db

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
)

// Wrapper is a wrapper over sql.DB
type Wrapper struct {
	db *pg.DB
}

// Comment represents the database version of a GitHub comment
type Comment struct {
	ID     int
	Data   string
	PullID int
	Pull   *Pull
}

// EComment represents the exported version of a GitHub comment
type EComment struct {
	ID               int    `json:"id"`
	Body             string `json:"body"`
	Position         int    `json:"position"`
	OriginalPosition int    `json:"original_position"`
	User             *EUser `json:"user"`
}

// Pull represents the database version of a GitHub pull request
type Pull struct {
	ID       int
	Data     string
	RepoID   int
	Repo     *Repo
	Comments []*Comment
}

// EPull represents the exported version of a Github pull request
type EPull struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Merged bool   `json:"merged"`
	User   *EUser `json:"user"`
	Repo   *ERepo `json:"repo"`
}

// EUser represents the exported version of a user in GitHub
type EUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// Repo represents the database version of a GitHub repository
type Repo struct {
	ID   int
	Data string
}

// ERepo represents the exported version of a GitHub repository
type ERepo struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	FullName      string  `json:"full_name"`
	Description   string  `json:"description"`
	WatchersCount int     `json:"watchers_count"`
	Language      string  `json:"language"`
	Owner         *EOwner `json:"owner"`
}

// EOwner represents the exported version of a GitHub repository
type EOwner struct {
	ID int `json:"id"`
}

// GetRepo is a function handler that retrieves a particular repository from the DB and writes it with the responseWriter
func (dbWrap *Wrapper) GetRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	repo := Repo{ID: ID}
	err = dbWrap.db.Select(&repo)
	if err != nil {
		panic(err)
	}

	var eRepo ERepo
	repoBytes := []byte(repo.Data)
	err = json.Unmarshal(repoBytes, &eRepo)
	if err != nil {
		panic(err)
	}

	repos := make([]ERepo, 1)
	repos[0] = eRepo

	response := buildRepoJSON(repos)

	addResponseHeaders(w)
	w.Write(response)
}

// GetRepos is a function handler that retrieves a set of repos from the DB and writes them with the responseWriter
func (dbWrap *Wrapper) GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []Repo
	err := dbWrap.db.Model(&repos).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	eRepos := make([]ERepo, len(repos))
	// Build JSON of the form {"repos": [...]}
	for idx, repo := range repos {
		var eRepo ERepo
		dataBytes := []byte(repo.Data)
		err = json.Unmarshal(dataBytes, &eRepo)
		if err != nil {
			panic(err)
		}
		eRepos[idx] = eRepo
	}

	response := buildRepoJSON(eRepos)

	addResponseHeaders(w)
	w.Write(response)
}

// GetPull is a function handler that retrieves a particular PR from the DB and writes it with the responseWriter
func (dbWrap *Wrapper) GetPull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["pullID"])
	if err != nil {
		panic(err)
	}

	var pull Pull
	err = dbWrap.db.Model(&pull).
		Column("pull.*", "Comments").
		Where("pull.id = ?", ID).
		Select()
	if err != nil {
		panic(err)
	}

	var ePull EPull
	pullBytes := []byte(pull.Data)
	err = json.Unmarshal(pullBytes, &ePull)
	if err != nil {
		panic(err)
	}

	ePulls := make([]EPull, 1)
	ePulls[0] = ePull

	comments := pull.Comments
	var eComments []EComment
	// Build JSON of the form {"comments": [...]}
	for idx, comment := range comments {
		var eComment EComment
		commentBytes := []byte(comment.Data)
		err = json.Unmarshal(commentBytes, &eComment)
		if err != nil {
			panic(err)
		}
		eComments[idx] = eComment
	}

	response := buildPullJSON(ePulls, eComments)

	addResponseHeaders(w)
	w.Write(response)
}

// GetPulls is a function handler that retrieves a set of PR's from the DB and writes them with the responseWriter
func (dbWrap *Wrapper) GetPulls(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	var pulls []Pull
	err = dbWrap.db.Model(&pulls).
		Where("pull.repo_id = ?", repoID).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	ePulls := make([]EPull, len(pulls))
	pullComments := pull.Comments
	// TODO: Find length of comment slice (all comments for all pulls)
	// eComments := make([]EComment, len(pullComments))

	for idx, pull := range pulls {
		// Build JSON of the form {"pulls": [...]}
		var ePull EPull
		dataBytes := []byte(pull.Data)
		err = json.Unmarshal(dataBytes, &ePull)
		if err != nil {
			panic(err)
		}
		pulls[idx] = pull

		// Build JSON of the form {"comments": [...]}
		for comment := range comments {
			var eComment EComment
			commentBytes := []byte(comment.Data)
			err = json.Unmarshal(commentBytes, &eComment)
			if err != nil {
				panic(err)
			}
			eComments[idx] = eComment
		}
	}

	response := buildPullJSON(ePulls, eComments)

	addResponseHeaders(w)
	w.Write(response)
}

// The builder functions have identical implementations except for the
// value passed into WrapModelJSON. Interfaces are more expensive
func buildRepoJSON(models []ERepo) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, model := range models {
		if &model != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			mJSON, err := json.Marshal(model)
			if err != nil {
				continue
			}
			buffer.Write(mJSON)
		}
	}
	buffer.WriteString(`]`)

	return wrapModelJSON("repos", buffer.Bytes())
}

func buildPullJSON(models []EPull, comments []EComment) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, model := range models {
		if &model != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			mJSON, err := json.Marshal(model)
			if err != nil {
				continue
			}
			buffer.Write(mJSON)
		}
	}
	buffer.WriteString(`]`)

	return wrapModelJSON("pulls", buffer.Bytes())
}

func buildCommentJSON(models []EComment) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, model := range models {
		if &model != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			mJSON, err := json.Marshal(model)
			if err != nil {
				continue
			}
			buffer.Write(mJSON)
		}
	}
	buffer.WriteString(`]`)

	return wrapModelJSON("comments", buffer.Bytes())
}

func wrapModelJSON(modelKey string, jsonBytes []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`{"`)
	buffer.WriteString(modelKey)
	buffer.WriteString(`":`)
	buffer.Write(jsonBytes)
	buffer.WriteString(`}`)

	return buffer.Bytes()
}

func addResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

// OpenDb initializes and returns a pointer to a Wrapper struct
func OpenDb() *Wrapper {
	user := os.Getenv("DEFAULT_POSTGRES_USER")
	db := pg.Connect(&pg.Options{
		User:     user,
		Database: "cleopatchra",
	})

	return &Wrapper{db: db}
}
