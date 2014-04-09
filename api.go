package addata

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hoisie/web"
)

// Queries submitted to the API must match the table name schema.
var (
	validQuery = regexp.MustCompile(`[a-zA-Z0-9-_]+`)
)

// API creates and managers a pointer to a web.Server, and has methods for the HTTP API so we can easily expose the other required types.
type API struct {
	Port   uint
	Server *web.Server
	Store  *Store
	Index  *Index
}

// NewAPI returns a pointer to a new API instance, creating all the required types for the API to function.
func NewAPI(port uint, i *Index, s *Store) *API {
	w := web.NewServer()

	w.Config.StaticDir = "static"

	a := &API{
		Port:   port,
		Server: w,
		Store:  s,
		Index:  i,
	}

	w.Get("/search", a.Search)
	w.Get("/table/([a-zA-Z0-9-_]+)", a.ReturnTable)

	return a
}

// Run starts the API instance, which in turn starts the web.Server.
func (a *API) Run() {
	log.Print("Starting API...")
	go a.Server.Run(fmt.Sprintf(":%v", a.Port))
}

// Search searches for a given table name against the Index, and returns a list of string tablenames.
func (a *API) Search(ctx *web.Context) []byte {
	term := ctx.Params["q"]

	if !validQuery.MatchString(term) {
		return []byte("[]")
	}

	q := Query{
		Term:       ctx.Params["q"],
		ReturnChan: make(chan []string),
	}

	a.Index.QueryChan <- q

	names := <-q.ReturnChan

	b, err := json.Marshal(names)

	if err != nil {
		ctx.Abort(500, "Could not build JSON list of names.")
		return nil
	}

	return b
}

// ReturnTable returns a CSV-encoded download of a table's contents.
func (a *API) ReturnTable(ctx *web.Context, tablename string) {
	results, err := a.Store.ReturnTable(tablename)

	if err != nil {
		ctx.Abort(500, err.Error())
	}

	ctx.SetHeader("Content-Description", "File Transfer", true)
	ctx.SetHeader("Content-Type", "text/csv; charset=utf-8", true)
	ctx.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%d.csv", tablename, time.Now().Unix()), true)

	w := csv.NewWriter(ctx)
	w.WriteAll(results)
	w.Flush()
}
