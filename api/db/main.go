package db

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/jfo84/cleopatchra/api/exports"
	"github.com/jfo84/cleopatchra/api/unmarshalling"
	"github.com/jinzhu/copier"
)

// Wrapper is a wrapper over pg.DB
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

// Pull represents the database version of a GitHub pull request
type Pull struct {
	ID       int
	Data     string
	RepoID   int
	Repo     *Repo
	Comments []*Comment
}

// Repo represents the database version of a GitHub repository
type Repo struct {
	ID   int
	Data string
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

	var eRepo exports.Repo
	repoBytes := []byte(repo.Data)
	err = json.Unmarshal(repoBytes, &eRepo)
	if err != nil {
		panic(err)
	}

	addResponseHeaders(w)

	if err := jsonapi.MarshalPayload(w, &eRepo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	// Build repo exports concurrently
	wg := &sync.WaitGroup{}
	eRepos := make([]*exports.Repo, len(repos))

	for idx, repo := range repos {
		wg.Add(1)

		go func(idx int, repo Repo) {
			defer wg.Done()

			var eRepo exports.Repo
			repoBytes := []byte(repo.Data)
			err = json.Unmarshal(repoBytes, &eRepo)
			if err != nil {
				panic(err)
			}

			eRepos[idx] = &eRepo
		}(idx, repo)
	}
	wg.Wait()

	addResponseHeaders(w)

	if err := jsonapi.MarshalPayload(w, eRepos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	var uPull unmarshalling.Pull
	pullBytes := []byte(pull.Data)
	err = json.Unmarshal(pullBytes, &uPull)
	if err != nil {
		panic(err)
	}

	var ePull exports.Pull
	copier.Copy(&ePull, &uPull)

	eComments := buildExportedComments(pull.Comments)
	ePull.Comments = eComments

	addResponseHeaders(w)

	if err := jsonapi.MarshalPayload(w, &ePull); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	// Build pull exports concurrently
	wg := &sync.WaitGroup{}
	ePulls := make([]*exports.Pull, len(pulls))

	for idx, pull := range pulls {
		wg.Add(1)

		go func(idx int, pull Pull) {
			defer wg.Done()

			var uPull unmarshalling.Pull
			pullBytes := []byte(pull.Data)
			err = json.Unmarshal(pullBytes, &uPull)
			if err != nil {
				panic(err)
			}

			var ePull exports.Pull
			copier.Copy(&ePull, &uPull)

			eComments := buildExportedComments(pull.Comments)
			ePull.Comments = eComments
			ePulls[idx] = &ePull
		}(idx, pull)
	}
	wg.Wait()

	addResponseHeaders(w)

	if err := jsonapi.MarshalPayload(w, ePulls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildExportedComments(comments []*Comment) []*exports.Comment {
	eComments := make([]*exports.Comment, len(comments))
	for idx, comment := range comments {
		var eComment exports.Comment
		commentBytes := []byte(comment.Data)
		err := json.Unmarshal(commentBytes, &eComment)
		if err != nil {
			panic(err)
		}
		eComments[idx] = &eComment
	}
	return eComments
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
