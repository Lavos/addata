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

var (
	validQuery = regexp.MustCompile(`[a-zA-Z0-9-_]+`)
)

type API struct {
	Port   uint
	Server *web.Server
	Store  *Store
	Index  *Index
}

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

func (a *API) Run() {
	log.Print("Starting API...")
	go a.Server.Run(fmt.Sprintf(":%v", a.Port))
}

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
