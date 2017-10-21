package db

import (
	"bytes"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// DB is a wrapper over sql.DB
type DB struct {
	db *sql.DB
}

// Pull represents a Github pull request
type Pull struct {
	id           int
	data, repoID *string
}

// Repo represents a Github repository
type Repo struct {
	id *string
}

func GetRepo(id int) *Repo {
	// TODO
}

func GetRepos(page int, perPage int) []*Repo {
	// TODO
}

func GetPull(id int) *Pull {
	db := openDb()
	rows, err := db.Query("SELECT * FROM pulls WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var (
		pullID 				int
		data, repoID	*string
	)

	rows.Next()
	err := rows.Scan(&id, &data, &repoID)
	if err != nil {
		panic(err)
	}

	p := &Pull{id:id, data:data, repoID:repoID}

	return p
}

func GetPulls(repoID *string, page int, perPage int) []*Pull {
	limit := perPage
	offset := page * perPage
	db := openDb()
	rows, err := db.Query("SELECT * FROM pulls WHERE repo_id = $1 LIMIT $2 OFFSET $3", repoID, limit, offset)
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
		err := rows.Scan(&id, &data)
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

func openDb() *sql.DB {
	connInfo := connectionInfo()
	
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	return db
}
