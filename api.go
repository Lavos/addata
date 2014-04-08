package addata

import (
	"log"
	"fmt"
	"time"
	"encoding/json"
	"encoding/csv"
	"regexp"

	"github.com/hoisie/web"
)

var (
	validQuery = regexp.MustCompile(`[a-zA-Z0-9-_]+`)
)


type API struct {
	Port uint
	server *web.Server
	store *Store
	index *Index
}

func newAPI(port uint, i *Index, s *Store) *API {
	w := web.NewServer()

	w.Config.StaticDir = "static"

	a := &API{
		Port: port,
		server: w,
		store: s,
		index: i,
	}

	w.Get("/search", a.search)
	w.Get("/table/([a-zA-Z0-9-_]+)", a.returnTable)

	return a
}

func (a *API) run() {
	log.Print("Starting API...")
	go a.server.Run(fmt.Sprintf(":%v", a.Port))
}

func (a *API) search(ctx *web.Context) []byte {
	term := ctx.Params["q"]

	if !validQuery.MatchString(term) {
		return []byte("[]")
	}

	q := Query{
		term: ctx.Params["q"],
		returnchan: make(chan []string),
	}

	a.index.querychan <- q

	names := <-q.returnchan

	b, err := json.Marshal(names);

	if err != nil {
		ctx.Abort(500, "Could not build JSON list of names.")
		return nil
	}

	return b
}

func (a *API) returnTable(ctx *web.Context, tablename string) {
	results, err := a.store.returnTable(tablename)

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
