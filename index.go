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

type Index struct {
	InvertedSuffix *ferret.InvertedSuffix

	QueryChan   chan Query
	RebuildChan chan []string
}

type Query struct {
	Term       string
	ReturnChan chan []string
}

func NewIndex() *Index {
	return &Index{
		QueryChan:   make(chan Query),
		RebuildChan: make(chan []string),
	}
}

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

func (i *Index) RebuildWith(names []string) {
	t := time.Now()
	dummy := make([]interface{}, len(names))

	i.InvertedSuffix = ferret.New(names, names, dummy, IndexConverter)
	log.Print("Created index in: ", time.Now().Sub(t))
}

func (i *Index) Query(term string) []string {
	t := time.Now()
	results, _ := i.InvertedSuffix.ErrorCorrectingQuery(term, 10, IndexCorrection)

	log.Print("Query completed in: ", time.Now().Sub(t))
	return results
}
