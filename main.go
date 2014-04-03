package main

import (
	"log"
	"os"
	"fmt"
	"time"
	"encoding/json"
	"encoding/csv"

	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/hoisie/web"
)

func setupDatabase () {
	conn, _ := sqlite3.Open("ads.db")
	defer conn.Close()

	conn.Exec("DROP TABLE IF EXISTS `apple-2014`")
	conn.Exec("CREATE TABLE `apple-2014` (id INTEGER PRIMARY KEY ASC, first_name TEXT, last_name TEXT, email TEXT)")

	conn.Exec("INSERT INTO `apple-2014` (first_name, last_name, email) VALUES ('Bob', 'Example,,,Tacos', 'bob@example.com')")
	conn.Commit()
}

func awaitQuitKey() {
	var buf [1]byte
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil || buf[0] == 'q' {
			return
		}
	}
}

func getList (ctx *web.Context) []byte {
	conn, _ := sqlite3.Open("ads.db")
	defer conn.Close()

	results := make([]string, 0)
	var name string
	row := make(sqlite3.RowMap)

	for s, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table'"); err == nil; err = s.Next() {
		s.Scan(row)
		name = row["name"].(string)
		results = append(results, name)
	}

	b, err := json.Marshal(results)

	if err != nil {
		ctx.Abort(500, "Could not get table data from database.")
		return []byte("")
	}

	return b
}

func returnTable (ctx *web.Context, tablename string) {
	conn, _ := sqlite3.Open("ads.db")
	defer conn.Close()

	s, err := conn.Query(fmt.Sprintf("SELECT * FROM `%s`", tablename))

	if err != nil {
		ctx.Abort(500, "Could not get data for the given table name.")
		return
	}

	ctx.SetHeader("Content-Description", "File Transfer", true)
	ctx.SetHeader("Content-Type", "text/csv; charset=utf-8", true)
	ctx.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%d.csv", tablename, time.Now().Unix()), true)

	w := csv.NewWriter(ctx)
	w.Write(s.Columns())

	row := make(sqlite3.RowMap)

	for err == nil {
		s.Scan(row)
		a := make([]string, 0, len(row))

		for _, value := range row {
			a = append(a, fmt.Sprint(value))
		}

		w.Write(a)
		err = s.Next()
	}

	w.Flush()
}

func main () {
	setupDatabase()

	w := web.NewServer()

	w.Get("/list", getList)
	w.Get("/table/([a-zA-Z0-9-_]+)", returnTable)

	log.Print("Starting server...")
	go w.Run(":8000")
	awaitQuitKey()
}
