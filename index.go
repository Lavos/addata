package addata

import (
	"time"
	"log"
	"github.com/argusdusty/Ferret"
)

var (
	IndexCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseLetters) }
	IndexSorter = func(s string, v interface{}, l int, i int) float64 { return -float64(l + i) }
	IndexConverter = func(s string) []byte { return []byte(s) }
)

type Index struct {
	InvertedSuffix *ferret.InvertedSuffix

	querychan chan Query
	rebuildchan chan []string
}

type Query struct {
	term string
	returnchan chan []string
}

func newIndex () *Index {
	return &Index{
		querychan: make(chan Query),
		rebuildchan: make(chan []string),
	}
}

func (i *Index) run() {
	for {
		select {
		case q := <-i.querychan:
			q.returnchan <- i.query(q.term)
		case n := <-i.rebuildchan:
			i.rebuildWith(n)
		}
	}
}

func (i *Index) rebuildWith(names []string) {
	t := time.Now()
	dummy := make([]interface{}, len(names))

        i.InvertedSuffix = ferret.New(names, names, dummy, IndexConverter)
	log.Print("Created index in: ", time.Now().Sub(t))
}

func (i *Index) query(term string) []string {
	t := time.Now()
	results, _ := i.InvertedSuffix.ErrorCorrectingQuery(term, 10, IndexCorrection)

	log.Print("Query completed in: ", time.Now().Sub(t))

	log.Printf("%#v - %v", results, len(results))

	return results
}
