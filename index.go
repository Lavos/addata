package addata

import (
	"github.com/argusdusty/Ferret"
	"log"
	"time"
)

var (
	IndexCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseLetters) }
	IndexSorter     = func(s string, v interface{}, l int, i int) float64 { return -float64(l + i) }
	IndexConverter  = func(s string) []byte { return []byte(s) }
)

// Index stores the internal ferret.InvertedSuffix and required channels for concurrent access.
type Index struct {
	InvertedSuffix *ferret.InvertedSuffix

	QueryChan   chan Query
	RebuildChan chan []string
}

// Query gathers the required information to make a query together for ease of communication across channels.
type Query struct {
	Term       string
	ReturnChan chan []string
}

// NewIndex returns a pointer to a new Index.
func NewIndex() *Index {
	return &Index{
		QueryChan:   make(chan Query),
		RebuildChan: make(chan []string),
	}
}

// Run starts the Index running it's required goroutines.
func (i *Index) Run() {
	go func() {
		for {
			select {
			case q := <-i.QueryChan:
				q.ReturnChan <- i.Query(q.Term)
			case n := <-i.RebuildChan:
				i.RebuildWith(n)
			}
		}
	}()
}

// RebuildWith rebuilds the Index's internal ferret.InvertedSuffix with the supplied slice of strings.
func (i *Index) RebuildWith(names []string) {
	t := time.Now()
	dummy := make([]interface{}, len(names))

	i.InvertedSuffix = ferret.New(names, names, dummy, IndexConverter)
	log.Print("Created index in: ", time.Now().Sub(t))
}

// Query checks a given string against the internal ferret.InvertedSuffix and return a slice of string matches.
func (i *Index) Query(term string) []string {
	t := time.Now()
	results, _ := i.InvertedSuffix.ErrorCorrectingQuery(term, 10, IndexCorrection)

	log.Print("Query completed in: ", time.Now().Sub(t))
	return results
}
