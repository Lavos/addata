package addata

import (
	"log"
	"fmt"
	"errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// dsn := "username:password@protocol(address)/dbname?param=value"
	dsn = "addata:faa1bcbea6a2e4d6@tcp(adbytes-db1.shared-prod.west1:3306)/addata"
)

type Store struct {
	dsn string
	db *sql.DB
}


func newStore(dsn string) *Store {
	db, err := sql.Open("mysql", dsn)
	log.Printf("%#v - %#v", db, err)

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return &Store{
		dsn: dsn,
		db: db,
	}
}

func (s *Store) getTableNames() []string {
	results := make([]string, 0)
	var name string

	rows, err := s.db.Query("SHOW TABLES")

	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&name)
		log.Print(name)
		results = append(results, name)
	}

	return results
}

func (s *Store) returnTable(tablename string) ([][]string, error) {
	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM `%s`", tablename))
	defer rows.Close()

	if err != nil {
		return nil, errors.New("Could not get data from the table name specified.")
	}

	columns, _ := rows.Columns()
	results := [][]string{columns}

	var (
		pointers []interface{}
		container []sql.RawBytes
		result []string
	)

	length := len(columns)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]sql.RawBytes, length)
		result = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)

		if err != nil {
			panic(err.Error())
		}

		for i, c := range container {
			if c == nil {
				result[i] = "-"
			} else {
				result[i] = string(c)
			}
		}

		results = append(results, result)
	}

	return results, nil
}
