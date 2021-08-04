package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Article struct {
	Id        int    `json:"id"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func main() {
	connURL := "postgres://postgres:password@0.0.0.0:5432/postgres?pool_max_conns=10"
	// connURL := "postgres://jokes:jokes@0.0.0.0:5432/jokes?pool_max_conns=10"
	dbpool, err := pgxpool.Connect(context.Background(), connURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	fmt.Println("Successfully connected!")
	var greeting string
	// err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed %v\n", err)
	// 	os.Exit(1)
	// }
	err = dbpool.QueryRow(context.Background(), "select datname from pg_database where datname='postgres'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
	dbpool.QueryRow(context.Background(), "CREATE TABLE IF NOT EXISTS dadjokes (id INT NOT NULL PRIMARY KEY, setup VARCHAR NOT NULL, punchline VARCHAR NOT NULL)").Scan()
	selectData, err := dbpool.Query(context.Background(), "SELECT * from dadjokes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "SELECT Query failed %v\n", err)
		os.Exit(1)
	}
	defer selectData.Close()
	var article Article

	for selectData.Next() {
		selectData.Scan(&article.Id, &article.Setup, &article.Punchline)
		fmt.Printf("%+v\n", article)
	}

	if selectData.Err() != nil {
		// if any error occurred while reading rows.
		fmt.Println("Error while reading the table: ", err)
		return
	}
	id := 3
	var sample Article
	// Select single element
	dbpool.QueryRow(context.Background(), "SELECT * from dadjokes WHERE id=$1", id).Scan(&sample.Id, &sample.Setup, &sample.Punchline)
	fmt.Printf("User with id=%v\n is %v\n %v\n", id, sample.Setup, sample.Punchline)

}

