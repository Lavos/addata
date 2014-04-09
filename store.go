package addata

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Type Store stores the DSN for the Database access, derived from the Application Configuration.
type Store struct {
	DSN string
	DB  *sql.DB
}

// NewStore returns a pointer to a new Store
func NewStore(dsn string) *Store {
	db, err := sql.Open("mysql", dsn)
	log.Printf("%#v - %#v", db, err)

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return &Store{
		DSN: dsn,
		DB:  db,
	}
}

// GetTableNames returns a list of the string table names, gathered from the database.
func (s *Store) GetTableNames() []string {
	results := make([]string, 0)
	var name string

	rows, err := s.DB.Query("SHOW TABLES")

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

// ReturnTable returns the rows of a table in a format that's applicable for CSV encoding in the API.
func (s *Store) ReturnTable(tablename string) ([][]string, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM `%s`", tablename))
	defer rows.Close()

	if err != nil {
		return nil, errors.New("Could not get data from the table name specified.")
	}

	columns, _ := rows.Columns()
	results := [][]string{columns}

	var (
		pointers  []interface{}
		container []sql.RawBytes
		result    []string
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
