package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable"
)

var testQuerries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB , err = sql.Open( dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuerries = New(testDB) // initate a new connection

	os.Exit(m.Run()) // run the test
}