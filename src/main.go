package main

import (
	"bytes"
	"database/sql"
	"os"
	"net/http"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// Pull represents a pull request
type Pull struct {
	id           int
	data, repoID string
}

func buildPulls(rows *sql.Rows) []*Pull {
	var (
		id           int
		data, repoID string
		pulls        []*Pull
	)

	for rows.Next() {
		i := 0
		err := rows.Scan(&id, &data, &repoID)
		if err != nil {
			panic(err)
		}

		p := &Pull{id:id, data:data, repoID: repoID}
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

func handleRepos(w http.ResponseWriter, r *http.Request) {
	
}

func handleRepo(w http.ResponseWriter, r *http.Request) {
	
}

func handlePulls(w http.ResponseWriter, r *http.Request) {
	
}

func handlePull(w http.ResponseWriter, r *http.Request) {

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func listenAndServe() {
	db := openDb()
	r := mux.NewRouter()

	r.HandleFunc("/repos", handleRepos)
	r.HandleFunc("/repo/{repoID}", handleRepo)
	r.HandleFunc("/repo/{repoID}/pulls", handlePulls)
	r.HandleFunc("/repo/{repoID}/pulls/{pullID}", handlePull)
	http.HandleFunc("/", handleIndex)

	addr := ":7000"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}

func openDb() *sql.DB {
	connInfo := connectionInfo()
	
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	listenAndServe()

	offset := 0
	rows, err := db.Query("SELECT * FROM pulls LIMIT 10 OFFSET $1", offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	pulls := buildPulls(rows)	
}
