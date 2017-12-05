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

// GetRepo is a function handler that retrieves a particular repository from the DB and writes it with the http.ResponseWriter
func (wrap *Wrapper) GetRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	repo := Repo{ID: ID}
	err = wrap.db.Select(&repo)
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

// GetRepos is a function handler that retrieves a set of repos from the DB and writes them with the http.ResponseWriter
func (wrap *Wrapper) GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []Repo
	err := wrap.db.Model(&repos).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	// Build repo exports concurrently
	wg := &sync.WaitGroup{}
	rLen := len(repos)
	eRepos := make([]*exports.Repo, rLen)
	wg.Add(rLen)

	for idx, repo := range repos {
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

// GetPull is a function handler that retrieves a particular PR from the DB and writes it with the http.ResponseWriter
func (wrap *Wrapper) GetPull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["pullID"])
	if err != nil {
		panic(err)
	}

	var pull Pull
	err = wrap.db.Model(&pull).
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

// GetPulls is a function handler that retrieves a set of PR's from the DB and writes them with the http.ResponseWriter
func (wrap *Wrapper) GetPulls(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	var pulls []Pull
	err = wrap.db.Model(&pulls).
		Where("pull.repo_id = ?", repoID).
		Apply(orm.Pagination(r.URL.Query())).
		Select()
	if err != nil {
		panic(err)
	}

	// Build pull exports concurrently
	wg := &sync.WaitGroup{}
	pLen := len(pulls)
	ePulls := make([]*exports.Pull, pLen)
	wg.Add(pLen)

	for idx, pull := range pulls {
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

// BeginTx begins a transaction on the wrapped pg.DB instance
func (wrap *Wrapper) BeginTx() *pg.Tx {
	tx, err := wrap.db.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

// OpenDB initializes and returns a pointer to a Wrapper struct
func OpenDB() *Wrapper {
	user := os.Getenv("DEFAULT_POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	// Assume no password means a password of empty string is fine
	password := os.Getenv("DEFAULT_POSTGRES_PASSWORD")
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: "cleopatchra",
	})

	return &Wrapper{db: db}
}

// OpenTestDB initializes and returns a pointer to a Wrapper struct and initializes temporary testing tables
func OpenTestDB() *Wrapper {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "cleopatchra_test",
	})
	wrapper := &Wrapper{db: db}

	err := createTestSchema(wrapper)
	if err != nil {
		panic(err)
	}

	return wrapper
}

func createTestSchema(wrapper *Wrapper) error {
	tables := []interface{}{
		&Repo{},
		&Pull{},
		&Comment{},
	}
	for _, table := range tables {
		err := wrapper.db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}

		err = wrapper.db.CreateTable(table, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
