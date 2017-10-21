package main

import (
	"bytes"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// Pull represents a pull request
type Pull struct {
	id           int
	data, repoID string
}

func listenAndServe(pulls []*Pull) {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "Not Found", 404)
	// 	return
	// })

	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(w, r)
	// })

	// addr := ":7000"
	// err := http.ListenAndServe(addr, nil)
	// if err != nil {
	// 	panic(err)
	// }
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

		p := &Pull{id: id, data: data, repoID: repoID}
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

func main() {
	connInfo := connectionInfo()

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	offset := 0
	rows, err := db.Query("SELECT * FROM pulls LIMIT 10 OFFSET $1", offset)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	pulls := buildPulls(rows)

	listenAndServe(pulls)
}
