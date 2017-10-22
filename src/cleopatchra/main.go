package main

import (
	"bytes"
	"database/sql"
	"os"
	"net/http"
	"db"
	"repo"
	"repos"
	"pulls"
	"pull"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func connectionInfo() string {
	var buffer bytes.Buffer

	buffer.WriteString("user=")
	user := os.Getenv("DEFAULT_POSTGRES_USER")
	buffer.WriteString(user)
	buffer.WriteString(" dbname=cleopatchra sslmode=disable")

	return buffer.String()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/repos", 301)
	return
}

func listenAndServe() {
	db := db.OpenDb()
	r := mux.NewRouter()

	reposController := repos.NewController(db)
	r.HandleFunc("/repos", reposController.Repos)

	repoController := repo.NewController(db)
	r.HandleFunc("/repo/{repoID}", repoController.Repo)

	pullsController := pulls.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls", pullsController.Pulls)

	pullController := pull.NewController(db)
	r.HandleFunc("/repo/{repoID}/pulls/{pullID}", pullController.Pull)

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
}
