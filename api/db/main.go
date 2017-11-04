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

// DataWrapper is a wrapper over string for passing through JSON string values during marshalling
type DataWrapper struct {
	Data string
}

// Comment represents a comment on a Github pull request
type Comment struct {
	ID     int          `json:"id"`
	Data   *DataWrapper `json:"data"`
	PullID int          `json:"pull_id"`
	Pull   *Pull        `json:"pull"`
}

// Pull represents a Github pull request
type Pull struct {
	ID       int          `json:"id"`
	Data     *DataWrapper `json:"data"`
	RepoID   int          `json:"repo_id"`
	Repo     *Repo        `json:"repo"`
	Comments []*Comment   `json:"comments"`
}

// Repo represents a Github repository
type Repo struct {
	ID   int          `json:"id"`
	Data *DataWrapper `json:"data"`
}

// NewComment is for initializing a Comment with a DataWrapper
func NewComment(ID int, Data string, PullID int) *Comment {
	return &Comment{
		ID:     ID,
		Data:   &DataWrapper{Data: Data},
		PullID: PullID,
	}
}

// NewPull is for initializing a Pull with a DataWrapper
func NewPull(ID int, Data string, RepoID int) *Pull {
	return &Pull{
		ID:     ID,
		Data:   &DataWrapper{Data: Data},
		RepoID: RepoID,
	}
}

// NewRepo is for initializing a Repo with a DataWrapper
func NewRepo(ID int, Data string) *Repo {
	return &Repo{
		ID:   ID,
		Data: &DataWrapper{Data: Data},
	}
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

	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []interface{} for wrapModelJSON
	models := make([]interface{}, 1)
	models[0] = &repo

	mJSON := buildModelJSON(models)
	response := wrapModelJSON("repos", mJSON)

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

	// Build JSON of the form {"repos": [...]}
	models := make([]interface{}, len(repos))
	for idx, repo := range repos {
		models[idx] = &repo
	}

	mJSON := buildModelJSON(models)
	response := wrapModelJSON("repos", mJSON)

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

	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []*string for buildModelJSON
	models := make([]interface{}, 1)
	models[0] = &pull

	mJSON := buildModelJSON(models)
	response := wrapModelJSON("pulls", mJSON)

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

	models := make([]interface{}, len(pulls))
	for idx, pull := range pulls {
		models[idx] = &pull
	}

	mJSON := buildModelJSON(models)
	response := wrapModelJSON("pulls", mJSON)

	addResponseHeaders(w)
	w.Write(response)
}

// MarshalJSON override to keep from re-marshalling the data JSON string
func (dWrap *DataWrapper) MarshalJSON() ([]byte, error) {
	return []byte(dWrap.Data), nil
}

func buildModelJSON(models []interface{}) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, model := range models {
		if model != nil {
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

	return buffer.Bytes()
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
