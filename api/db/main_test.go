package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/bluele/factory-go/factory"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var PullFactory = factory.NewFactory(
	&Pull{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	pull := args.Instance().(*Pull)
	return fmt.Sprintf("pull-%d", pull.ID), nil
}).OnCreate(func(args factory.Args) error {
	tx := args.Context().Value("tx").(*pg.Tx)
	return tx.Insert(args.Instance())
}).SubFactory("Repo", RepoFactory)

var RepoFactory = factory.NewFactory(
	&Repo{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Data", func(args factory.Args) (interface{}, error) {
	repo := args.Instance().(*Repo)
	return fmt.Sprintf("repo-%d", repo.ID), nil
}).OnCreate(func(args factory.Args) error {
	tx := args.Context().Value("tx").(*pg.Tx)
	return tx.Insert(args.Instance())
})

func createTestSchema(db *pg.DB) error {
	tables := []interface{}{
		&Repo{},
		&Pull{},
	}
	for _, table := range tables {
		err := db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}

		err = db.CreateTable(table, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func openDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "cleopatchra_test",
	})

	err := createTestSchema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func TestDB(t *testing.T) {
	db := openDB()
	for i := 0; i < 3; i++ {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		// To avoid collisions with other keys
		type key string
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
