package db

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Wrapper is a wrapper over sql.DB
type Wrapper struct {
	db *pg.DB
}

// TODO: Combine these types?? Much of the code for iterating through pulls/repos
// could be generalized if this was done. Feels too early to do so now

// Pull represents a Github pull request
type Pull struct {
	id   int
	data string
}

// Repo represents a Github repository
type Repo struct {
	id   int
	data string
}

// GetRepo is a function handler that retrieves a particular repository from the DB and writes it with the responseWriter
func (dbWrap *Wrapper) GetRepo(w http.ResponseWriter, r *http.Request) {
	repoIDs, ok := r.URL.Query()["repoID"]
	if !ok || len(repoIDs) < 1 {
		panic("No repoID in repos query")
	}
	id, err := strconv.Atoi(repoIDs[0])
	if err != nil {
		panic(err)
	}

	repo := Repo{id: id}
	err = dbWrap.db.Select(&repo)
	if err != nil {
		panic(err)
	}

	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []*string for buildModelJSON
	mStrings := make([]*string, 1)
	mStrings[0] = &repo.data

	mJSON := buildModelJSON(mStrings)
	response := wrapModelJSON("repos", mJSON)

	addResponseHeaders(w)
	w.Write(response)
}

// GetRepos is a function handler that retrieves a set of repos from the DB and writes them with the responseWriter
func (dbWrap *Wrapper) GetRepos(w http.ResponseWriter, r *http.Request) {
	var repos []Repo
	// query := dbWrap.db.Model(&repos).Apply(orm.Pagination(r.URL.Query()))
	query := orm.NewQuery(nil, &Repo{})
	query = query.Apply(orm.Pagination(r.URL.Query()))
	err := query.Select()
	if err != nil {
		panic(err)
	}

	// Build JSON of the form {"repos": [...]}
	fmt.Println(r.URL.Query())
	fmt.Println(query.getFields())
	mStrings := make([]*string, len(repos))
	for idx, repo := range repos {
		fmt.Printf("%s\n", strconv.Itoa(repo.id))
		fmt.Printf("%s\n", repo.data)
		mStrings[idx] = &repo.data
	}

	mJSON := buildModelJSON(mStrings)
	response := wrapModelJSON("repos", mJSON)

	addResponseHeaders(w)
	w.Write(response)
}

// GetPull is a function handler that retrieves a particular PR from the DB and writes it with the responseWriter
func (dbWrap *Wrapper) GetPull(w http.ResponseWriter, r *http.Request) {
	pullIDs, ok := r.URL.Query()["pullID"]
	if !ok || len(pullIDs) < 1 {
		panic("No pullID in pulls query")
	}
	id, err := strconv.Atoi(pullIDs[0])
	if err != nil {
		panic(err)
	}

	pull := Pull{id: id}
	err = dbWrap.db.Select(&pull)
	if err != nil {
		panic(err)
	}

	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []*string for buildModelJSON
	mStrings := make([]*string, 1)
	mStrings[0] = &pull.data

	mJSON := buildModelJSON(mStrings)
	response := wrapModelJSON("pulls", mJSON)

	addResponseHeaders(w)
	w.Write(response)
}

// GetPulls is a function handler that retrieves a set of PR's from the DB and writes them with the responseWriter
func (dbWrap *Wrapper) GetPulls(w http.ResponseWriter, r *http.Request) {
	var pulls []Pull
	err := dbWrap.db.Model(&pulls).Apply(orm.Pagination(r.URL.Query())).Select()
	if err != nil {
		panic(err)
	}

	// Build JSON of the form {"pulls": [...]}
	mStrings := make([]*string, len(pulls))
	for idx, pull := range pulls {
		mStrings[idx] = &pull.data
	}

	mJSON := buildModelJSON(mStrings)
	response := wrapModelJSON("pulls", mJSON)

	addResponseHeaders(w)
	w.Write(response)
}

func buildModelJSON(modelStrings []*string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, modelString := range modelStrings {
		if modelString != nil {
			if idx != 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(*modelString)
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
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "",
		Database: "cleopatchra",
	})

	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	return &Wrapper{db: db}
}
