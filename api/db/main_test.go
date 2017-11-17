package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	"github.com/jfo84/factory-go/factory"
)

// To avoid collisions with other keys
type key string

// PullFactory is a factory for generating temporary rows on the pulls table
var PullFactory = factory.NewFactory(
	&Pull{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	pull := args.Instance().(*Pull)
	return fmt.Sprintf("pull-%d", pull.ID), nil
}).OnCreate(func(args factory.Args) error {
	const txKey key = "tx"
	tx := args.Context().Value(txKey).(*pg.Tx)
	return tx.Insert(args.Instance())
}).SubFactory("Repo", RepoFactory)

// RepoFactory is a factory for generating temporary rows on the repos table
var RepoFactory = factory.NewFactory(
	&Repo{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	repo := args.Instance().(*Repo)
	return fmt.Sprintf("repo-%d", repo.ID), nil
}).OnCreate(func(args factory.Args) error {
	const txKey key = "tx"
	tx := args.Context().Value(txKey).(*pg.Tx)
	return tx.Insert(args.Instance())
})

func TestDB(t *testing.T) {
	dbWrap := openTestDB()

	for i := 0; i < 3; i++ {
		tx, err := dbWrap.db.Begin()
		if err != nil {
			panic(err)
		}

		const txKey key = "tx"

		ctx := context.WithValue(context.Background(), txKey, tx)
		v, err := PullFactory.CreateWithContext(ctx)
		if err != nil {
			panic(err)
		}
		pull := v.(*Pull)
		fmt.Println(pull, *pull.Repo)
		tx.Commit()
	}
}
