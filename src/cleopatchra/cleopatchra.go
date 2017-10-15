package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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

		idString := strconv.Itoa(id)

		fmt.Printf("%v\n", idString)

		type Pull struct {
			number int
			url    string
		}
		dec := json.NewDecoder(strings.NewReader(data))
		for {
			var p Pull
			if err := dec.Decode(&p); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			fmt.Printf("URL: %s\n", p.url)
		}
	}
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

	payload := buildPayload(rows)
}
