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

type comment struct {
	ID     int
	Data   string
	pullID int
	pull   *pull
}

// Comment represents the sliced version of a Github comment
type Comment struct {
	id               int
	body             string
	position         int
	originalPosition int
	user             *User
}

type pull struct {
	ID       int
	Data     string
	repoID   int
	repo     *repo
	comments []*comment
}

// Pull represents the sliced version of a Github pull request
type Pull struct {
	id       int
	number   int
	title    string
	body     string
	merged   bool
	user     *User
	repo     *Repo
	comments *Comment
}

// User represents a user in GitHub
type User struct {
	id    int
	login string
}

type repo struct {
	ID   int
	Data string
}

// Repo represents the sliced version of a Github repository
type Repo struct {
	id            int
	name          string
	fullName      string
	description   string
	watchersCount int
	language      string
	owner         *Owner
}

// Owner represents the owner of a Repo
type Owner struct {
	id int
}

// GetRepo is a function handler that retrieves a particular repository from the DB and writes it with the responseWriter
func (dbWrap *Wrapper) GetRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	repoDB := repo{ID: ID}
	err = dbWrap.db.Select(&repoDB)
	if err != nil {
		panic(err)
	}

	var repo *Repo
	dataBytes := []byte(repoDB.Data)
	err = json.Unmarshal(dataBytes, &repo)
	if err != nil {
		panic(err)
	}

	repos := make([]*Repo, 1)
	repos[0] = repo

	response := buildRepoJSON(repos)

	addResponseHeaders(w)
	w.Write(response)
}

// GetRepos is a function handler that retrieves a set of repos from the DB and writes them with the responseWriter
func (dbWrap *Wrapper) GetRepos(w http.ResponseWriter, r *http.Request) {
	var repoDBs []repo
	err := dbWrap.db.Model(&repoDBs).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	repos := make([]*Repo, len(repoDBs))
	// Build JSON of the form {"repos": [...]}
	for idx, repoDB := range repoDBs {
		var repo *Repo
		dataBytes := []byte(repoDB.Data)
		err = json.Unmarshal(dataBytes, &repo)
		if err != nil {
			panic(err)
		}
		repos[idx] = repo
	}

	response := buildRepoJSON(repos)

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

	var pullDB pull
	err = dbWrap.db.Model(&pullDB).
		Column("pull.*", "Comments").
		Where("pull.id = ?", ID).
		Select()
	if err != nil {
		panic(err)
	}

	var pull *Pull
	dataBytes := []byte(pullDB.Data)
	err = json.Unmarshal(dataBytes, &pull)
	if err != nil {
		panic(err)
	}

	pulls := make([]*Pull, 1)
	pulls[0] = pull

	response := buildPullJSON(pulls)

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

	var pullDBs []pull
	err = dbWrap.db.Model(&pullDBs).
		Where("pull.repo_id = ?", repoID).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	pulls := make([]*Pull, len(pullDBs))
	// Build JSON of the form {"pulls": [...]}
	for idx, pullDB := range pullDBs {
		var pull *Pull
		dataBytes := []byte(pullDB.Data)
		err = json.Unmarshal(dataBytes, &pull)
		if err != nil {
			panic(err)
		}
		pulls[idx] = pull
	}

	response := buildPullJSON(pulls)

	addResponseHeaders(w)
	w.Write(response)
}

func buildRepoJSON(repos []*Repo) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, repo := range repos {
		if repo != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			rJSON, err := json.Marshal(repo)
			if err != nil {
				continue
			}
			buffer.Write(rJSON)
		}
	}
	buffer.WriteString(`]`)

	repoBytes := wrapModelJSON("repos", buffer.Bytes())

	return repoBytes
}

func buildPullJSON(pulls []*Pull) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, pull := range pulls {
		if pull != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			pJSON, err := json.Marshal(pull)
			if err != nil {
				continue
			}
			buffer.Write(pJSON)
		}
	}
	buffer.WriteString(`]`)

	pullBytes := wrapModelJSON("pulls", buffer.Bytes())

	return pullBytes
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
