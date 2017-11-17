package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg"
	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/pull"
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
	return fmt.Sprintf("pull-%d", pull.ID), nil
}).OnCreate(func(args factory.Args) error {
	const txKey key = "tx"
	tx := args.Context().Value(txKey).(*pg.Tx)
	return tx.Insert(args.Instance())
}).SubSliceFactory("Comments", CommentFactory, func() int { return 3 })

// CommentFactory is a factory for generating temporary rows on the comments table
var CommentFactory = factory.NewFactory(
	&db.Pull{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	comment := args.Instance().(*db.Comment)
	return fmt.Sprintf("comment-%d", comment.ID), nil
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
		fmt.Println(pull, pull.Comments)
		tx.Commit()
	}

	req, err := http.NewRequest("GET", "/pulls/1", nil)

	checkError(err, t)

	rr := httptest.NewRecorder()

	pullController := pull.NewController(dbWrap)
	http.HandlerFunc(pullController.Get).ServeHTTP(rr, req)

	// Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	// Confirm the returned json is what we expected
	// Manually build up the expected json string
	expected := string(`[{"id":1,"title":"New blog resolution","content":"I have decided to give my blog a new life and would hence forth try to write as often"},{"id":2,"title":"Go is cool","content":"Yeah i have been told that multiple times"},{"id":3,"title":"Interminttent fasting","content":"You should try this out, it helps clear the brain and tons of health benefits"},{"id":4,"title":"Yet another blog post","content":"I made a resolution earlier to keep on writing. Here is an affirmation of that"},{"id":5,"title":"Backpacking","content":"Yup, i did just that"}]`)

	// The assert package checks if both JSON string are equal and for a plus, it actually confirms if our manually built JSON string is valid
	assert.JSONEq(t, expected, rr.Body.String(), "Response body differs")
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}
