package db

import (
	"bytes"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// DBWrapper is a wrapper over sql.DB
type DBWrapper struct {
	db *sql.DB
}

// Pull represents a Github pull request
type Pull struct {
	id           int
	data, repoID *string
}

// Repo represents a Github repository
type Repo struct {
	id		int
	data	*string
}

func (dbwrap *DBWrapper) GetRepo(id int) *Repo {
	dbwrap = OpenDb()
	rows, err := dbwrap.db.Query("SELECT * FROM repos WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var data *string

	rows.Next()
	err = rows.Scan(&id, data)
	if err != nil {
		panic(err)
	}

	r := &Repo{id:id, data:data}

	return r
}

func (dbwrap *DBWrapper) GetRepos(page int, perPage int) []*Repo {
	limit := perPage
	offset := page * perPage
	dbwrap = OpenDb()
	rows, err := dbwrap.db.Query("SELECT * FROM repos LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var (
		id           int
		data 				 *string
		repos        []*Repo
	)

	for rows.Next() {
		i := 0
		err := rows.Scan(&id, data)
		if err != nil {
			panic(err)
		}

	r := &Repo{id:id, data:data}
		repos[i] = r
		i++
	}

	return repos
}

func (dbwrap *DBWrapper) GetPull(id int) *Pull {
	dbwrap = OpenDb()
	rows, err := dbwrap.db.Query("SELECT * FROM pulls WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var data, repoID	*string

	rows.Next()
	err = rows.Scan(&id, data, repoID)
	if err != nil {
		panic(err)
	}

	p := &Pull{id:id, data:data, repoID:repoID}

	return p
}

func (dbwrap *DBWrapper) GetPulls(repoID *string, page int, perPage int) []*Pull {
	limit := perPage
	offset := page * perPage
	dbwrap = OpenDb()
	rows, err := dbwrap.db.Query("SELECT * FROM pulls WHERE repo_id = $1 LIMIT $2 OFFSET $3", repoID, limit, offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var (
		id           int
		data 				 *string
		pulls        []*Pull
	)

	for rows.Next() {
		i := 0
		err := rows.Scan(&id, data)
		if err != nil {
			panic(err)
		}

		p := &Pull{id:id, data:data, repoID:repoID}
		pulls[i] = p
		i++
	}

	return pulls
}

func connectionInfo() string {
	var buffer bytes.Buffer

	buffer.WriteString("user=")
	user := os.Getenv("DEFAULT_POSTGRES_USER")
	buffer.WriteString(user)
	buffer.WriteString(" dbname=cleopatchra sslmode=disable")

	return buffer.String()
}

func OpenDb() *DBWrapper {
	connInfo := connectionInfo()
	
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	return &DBWrapper{db:db}
}
