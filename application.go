package addata

import (
	"fmt"
	"time"
)

// Application represents a collection of the other types within this package.
// This is the only type that most implementations will have to create, usually via NewApplication.
type Application struct {
	API   *API
	Store *Store
	Index *Index
}

// Configuration maps directly to the passed JSON configuration file keys.
type Configuration struct {
	Username, Password, Hostname, Database string
	DBPort, APIPort                        uint
}

// NewApplication returns a pointer to a new Application.
// Most implemenations will only need to call this to create all the required other types.
func NewApplication(c *Configuration) *Application {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.Username, c.Password, c.Hostname, c.DBPort, c.Database)

	s := NewStore(dsn)
	i := NewIndex()
	a := NewAPI(c.APIPort, i, s)

	app := &Application{
		API:   a,
		Store: s,
		Index: i,
	}

	return app
}

// Maintenance creates a loop that refreshs the index on an interval.
func (a *Application) Maintenance() {
	c := time.Tick(10 * time.Second)

	go func() {
		for {
			a.Index.RebuildChan <- a.Store.GetTableNames()
			<-c
		}
	}()
}

// Run starts the Application, which then starts all the other required types within.
func (a *Application) Run() {
	a.API.Run()
	a.Index.Run()
	a.Maintenance()
}
