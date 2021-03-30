package main

import (
	"database/sql"
	"fmt"
	"github.com/alextsa22/gophercises/08-phone/db"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "gophercises_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	opendb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer opendb.Close()

	psqldb, err := db.NewPostgresDB(opendb, dbname)
	if err != nil {
		log.Fatal(err)
	}

	if err := psqldb.InitExample(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("before normalize:")
	ShowAllRows(psqldb)

	if err := psqldb.Normalize(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("after normalize:")
	ShowAllRows(psqldb)
}

func ShowAllRows(db *db.PostgresDB) {
	phones, err := db.AllPhones()
	if err != nil {
		log.Fatal()
	}

	for _, p := range phones {
		fmt.Println(p)
	}
}
