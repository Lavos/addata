package addata

import (
	"time"
	"fmt"
)

type Application struct {
	api *API
	store *Store
	index *Index
}

type Configuration struct {
	Username, Password, Hostname, Database string
	DBPort, APIPort uint
}

func NewApplication (c *Configuration) *Application {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.Username, c.Password, c.Hostname, c.DBPort, c.Database)

	s := newStore(dsn)
	i := newIndex()
	a := newAPI(c.APIPort, i, s)

	app := &Application{
		api: a,
		store: s,
		index: i,
	}

	return app
}

func (a *Application) Maintenance () {
	c := time.Tick(10 * time.Second)

	for {
		<-c
		a.index.rebuildchan <- a.store.getTableNames()
	}
}

func (a *Application) Run() {
	a.api.run()
	go a.index.run()
	a.index.rebuildchan <- a.store.getTableNames()
	go a.Maintenance()
}
