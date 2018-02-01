package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/pull"
	"github.com/jfo84/cleopatchra/api/pulls"
	"github.com/jfo84/factory-go/factory"
	"github.com/stretchr/testify/assert"
)

// To avoid collisions with other keys
type key string

// PullFactory is a factory for generating temporary rows on the pulls table
var PullFactory = factory.NewFactory(
	&db.Pull{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	pull := args.Instance().(*db.Pull)
	fileName := fmt.Sprintf("./testing/fixtures/pulls/%d.json", pull.ID)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data[:]), nil
}).OnCreate(func(args factory.Args) error {
	const txKey key = "tx"
	tx := args.Context().Value(txKey).(*pg.Tx)
	return tx.Insert(args.Instance())
}).SubSliceFactory("Comments", CommentFactory, func() int { return 2 })

// CommentFactory is a factory for generating temporary rows on the comments table
var CommentFactory = factory.NewFactory(
	&db.Comment{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).SeqInt("PullID", func(n int) (interface{}, error) {
	// Return 1, 1, 2, 2, 3, 3, etc.
	// TODO: Maybe pass around context or add a way to access parent factory values
	if n%2 == 1 {
		return (n + 1) / 2, nil
	}
	return n / 2, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	comment := args.Instance().(*db.Comment)
	fileName := fmt.Sprintf("./testing/fixtures/comments/%d.json", comment.ID)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data[:]), nil
}).OnCreate(func(args factory.Args) error {
	const txKey key = "tx"
	tx := args.Context().Value(txKey).(*pg.Tx)
	return tx.Insert(args.Instance())
})

func TestCleopatchra(t *testing.T) {
	dbWrap := db.OpenTestDB()

	for i := 0; i < 3; i++ {
		tx := dbWrap.BeginTx()

		const txKey key = "tx"

		ctx := context.WithValue(context.Background(), txKey, tx)
		v, err := PullFactory.CreateWithContext(ctx)
		if err != nil {
			panic(err)
		}
		pull := v.(*db.Pull)
		fmt.Printf("%d", pull.ID)
		tx.Commit()
	}

	router := mux.NewRouter().StrictSlash(true)

	req, err := http.NewRequest("GET", "/pulls/1", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	recorder := httptest.NewRecorder()

	pullController := pull.NewController(dbWrap)
	router.HandleFunc("/pulls/{pullID}", pullController.Get)
	router.ServeHTTP(recorder, req)

	// Confirm the response has the right status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	// Confirm the returned json is what we expected
	eBytes, err := ioutil.ReadFile("./testing/fixtures/expected_pull.json")
	if err != nil {
		panic(err)
	}
	expected := string(eBytes[:])

	fmt.Println(recorder.Body.String())

	assert.JSONEqf(t, expected, recorder.Body.String(), "Response body differs")

	req, err = http.NewRequest("GET", "/repos/1/pulls", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	recorder = httptest.NewRecorder()

	pullsController := pulls.NewController(dbWrap)
	router.HandleFunc("/repos/{repoID}/pulls", pullsController.Get)
	router.ServeHTTP(recorder, req)

	// Confirm the response has the right status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	// Confirm the returned json is what we expected
	eBytes, err = ioutil.ReadFile("./testing/fixtures/expected_pulls.json")
	if err != nil {
		panic(err)
	}
	expected = string(eBytes[:])

	fmt.Println(recorder.Body.String())

	assert.JSONEqf(t, expected, recorder.Body.String(), "Response body differs")
}
