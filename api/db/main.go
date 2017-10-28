package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	// Adds pq bindings to database/sql
	_ "github.com/lib/pq"
)

// Wrapper is a wrapper over sql.DB
type Wrapper struct {
	db *sql.DB
}

// TODO: Combine these types?? Much of the code for iterating through pulls/repos
// could be generalized if this was done. Feels too early to do so now

// Pull represents a Github pull request
type Pull struct {
	id   int
	data *string
}

// Repo represents a Github repository
type Repo struct {
	id   int
	data *string
}

// GetRepo is a function handler that retrieves a particular repository from the DB,
// marshalls it to JSON, and writes it with the responseWriter
func (dbWrap *Wrapper) GetRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["repoID"])
	if err != nil {
		panic(err)
	}

	rows, err := dbWrap.db.Query("SELECT * FROM repos WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var data string

	rows.Next()
	err = rows.Scan(&id, &data)
	if err != nil {
		panic(err)
	}
	if id == 0 {
		return
	}

	repo := &Repo{id: id, data: &data}
	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []*string for buildModelJSON
	repoStrings := make([]*string, 1)
	repoStrings[0] = repo.data

	mJSON := buildModelJSON(repoStrings)
	response := wrapModelJSON("repos", mJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetRepos is a function handler that retrieves a set of repos from the DB,
// marshalls them to JSON, and writes them with the responseWriter
func (dbWrap *Wrapper) GetRepos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Apply defaults of page 1 and perPage 10
	var (
		page, perPage int
		err           error
	)

	if vars["page"] != "" {
		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			panic(err)
		}
	} else {
		page = 1
	}

	if vars["perPage"] != "" {
		perPage, err = strconv.Atoi(vars["perPage"])
		if err != nil {
			panic(err)
		}
	} else {
		perPage = 10
	}

	limit := perPage
	offset := (page * perPage) - perPage

	rows, err := dbWrap.db.Query("SELECT * FROM repos LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var (
		id   int
		data string
	)

	// Build JSON of the form {"repos": [...]}
	repos := make([]*string, perPage)

	i := 0
	for rows.Next() {
		err := rows.Scan(&id, &data)
		if err != nil {
			panic(err)
		}
		if id == 0 {
			return
		}

		repo := &Repo{id: id, data: &data}

		repos[i] = repo.data
		i++
	}
	mJSON := buildModelJSON(repos)
	response := wrapModelJSON("repos", mJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetPull is a function handler that retrieves a particular pull request from the DB,
// marshalls it to JSON, and writes it with the responseWriter
func (dbWrap *Wrapper) GetPull(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["pullID"])
	if err != nil {
		panic(err)
	}

	rows, err := dbWrap.db.Query("SELECT * FROM pulls WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var data, repoID string

	rows.Next()
	err = rows.Scan(&id, &data, &repoID)
	if err != nil {
		panic(err)
	}
	if id == 0 {
		return
	}

	p := &Pull{id: id, data: &parsedData}
	// In order to keep the builder interface agnostic, I need to
	// generate a one-dimensional []*string for buildModelJSON
	pullStrings := make([]*string, 1)
	pullStrings[0] = p.data

	mJSON := buildModelJSON(pullStrings)
	response := wrapModelJSON("pulls", mJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetPulls is a function handler that retrieves a set of pull requests from the DB,
// marshalls them to JSON, and writes them with the responseWriter
func (dbWrap *Wrapper) GetPulls(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Apply defaults of page 1, perPage 10, and repoID "facebook/react"
	var (
		page, perPage int
		repoID        string
		err           error
	)

	if vars["page"] != "" {
		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			panic(err)
		}
	} else {
		page = 1
	}

	if vars["perPage"] != "" {
		perPage, err = strconv.Atoi(vars["perPage"])
		if err != nil {
			panic(err)
		}
	} else {
		perPage = 10
	}

	if vars["repoID"] != "" {
		repoID = vars["repoID"]
		if err != nil {
			panic(err)
		}
	} else {
		// TODO: Remove default
		repoID = "facebook/react"
	}

	limit := perPage
	offset := (page * perPage) - perPage

	rows, err := dbWrap.db.Query("SELECT * FROM pulls WHERE repo_id = $1 LIMIT $2 OFFSET $3", repoID, limit, offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var (
		id   int
		data string
	)

	// Build JSON of the form {"pulls": [...]}
	pulls := make([]*string, perPage)

	i := 0
	for rows.Next() {
		err := rows.Scan(&id, &data, &repoID)
		if err != nil {
			panic(err)
		}
		if id == 0 {
			return
		}

		parsedData, err := strconv.Unquote(data)
		if err != nil {
			panic(err)
		}

		p := &Pull{id: id, data: &parsedData}
		pulls[i] = p.data

		i++
	}

	mJSON := buildModelJSON(pulls)
	response := wrapModelJSON("repos", mJSON)

	mJSON, err := json.Marshal(pulls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := wrapModelJSON("pulls", mJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// TODO: Remove these. The only struct-based solution I could find
// would require unmarshalling and then marshalling back to JSON
func buildModelJSON(modelStrings []*string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for idx, string := range modelStrings {
		if idx != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(*string)
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

func connectionInfo() string {
	var buffer bytes.Buffer

	buffer.WriteString("user=")
	user := os.Getenv("DEFAULT_POSTGRES_USER")
	buffer.WriteString(user)
	buffer.WriteString(" dbname=cleopatchra sslmode=disable")

	return buffer.String()
}

// OpenDb initializes and returns a pointer to a Wrapper struct
func OpenDb() *Wrapper {
	connInfo := connectionInfo()

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	return &Wrapper{db: db}
}
