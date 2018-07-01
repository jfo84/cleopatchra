package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/go-pg/pg"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/jfo84/cleopatchra/api/db"
	"github.com/jfo84/cleopatchra/api/exports"
	"github.com/jfo84/cleopatchra/api/pull"
	"github.com/jfo84/cleopatchra/api/pulls"
	"github.com/jfo84/factory-go/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
}).SeqInt("RepoID", func(n int) (interface{}, error) {
	// Return 1, 1, 2, 2, 3, 3, etc.
	// TODO: Maybe pass around context or add a way to access parent factory values
	if n%2 == 1 {
		return (n + 1) / 2, nil
	}
	return n / 2, nil
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

var _ = Describe("TestCleopatchra", func() {
	var (
		req *http.Request
		err error
	)

	router := mux.NewRouter().StrictSlash(true)
	dbWrap := db.OpenTestDB()

	for i := 0; i < 3; i++ {
		tx := dbWrap.BeginTx()

		const txKey key = "tx"

		ctx := context.WithValue(context.Background(), txKey, tx)
		_, err := PullFactory.CreateWithContext(ctx)
		if err != nil {
			panic(err)
		}

		tx.Commit()
	}

	Context("Pull Requests", func() {
		It("Should correctly return a pull", func() {
			req, err = http.NewRequest("GET", "/pulls/1", nil)

			recorder := httptest.NewRecorder()

			pullController := pull.NewController(dbWrap)
			router.HandleFunc("/pulls/{pullID}", pullController.Get)
			router.ServeHTTP(recorder, req)

			// Confirm the returned json is what we expected
			var eBytes []byte
			eBytes, err = ioutil.ReadFile("./testing/fixtures/expected_pull.json")
			if err != nil {
				panic(err)
			}

			expectedPull := new(exports.Pull)
			reader := bytes.NewReader(eBytes)

			jsonapi.UnmarshalPayload(reader, expectedPull)

			reader = bytes.NewReader(recorder.Body.Bytes())
			actualPull := new(exports.Pull)

			jsonapi.UnmarshalPayload(reader, actualPull)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(expectedPull).To(BeEquivalentTo(actualPull))
		})

		It("Should correctly return pulls", func() {
			req, err = http.NewRequest("GET", "/repos/1/pulls", nil)

			recorder := httptest.NewRecorder()

			pullsController := pulls.NewController(dbWrap)
			router.HandleFunc("/repos/{repoID}/pulls", pullsController.Get)
			router.ServeHTTP(recorder, req)

			// Confirm the returned json is what we expected
			eBytes, err := ioutil.ReadFile("./testing/fixtures/expected_pulls.json")
			if err != nil {
				panic(err)
			}

			var ePulls []*exports.Pull
			pullType := reflect.TypeOf(ePulls)

			reader := bytes.NewReader(eBytes)

			expectedPulls, err := jsonapi.UnmarshalManyPayload(reader, pullType)

			reader = bytes.NewReader(recorder.Body.Bytes())

			actualPulls, err := jsonapi.UnmarshalManyPayload(reader, pullType)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(expectedPulls).To(BeEquivalentTo(actualPulls))
		})
	})
})
