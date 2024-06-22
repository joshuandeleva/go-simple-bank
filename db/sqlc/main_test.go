package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joshuandeleva/simplebank/util"
	_ "github.com/lib/pq"
)

var testQuerries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config , err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	testDB , err = sql.Open( config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuerries = New(testDB) // initate a new connection

	os.Exit(m.Run()) // run the test
}