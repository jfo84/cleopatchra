package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func buildPayload(rows *sql.Rows) {
	var (
		id           int
		data, repoID string
	)

	for rows.Next() {
		err := rows.Scan(&id, &data, &repoID)
		if err != nil {
			panic(err)
		}

		type Pull struct {
			number                                  int
			url, state, title, body, mergeCommitSha string
		}
		dec := json.NewDecoder(strings.NewReader(data))
		for {
			var p Pull
			if err := dec.Decode(&p); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if p.number == 10799 {
				fmt.Printf("%v\n", data)
			}

			fmt.Printf("Merge Commit SHA: %s\n", p.title)
		}
	}
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

	buildPayload(rows)
}
