package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jw-dev/x/chessquery/pkg/pgn"
	"github.com/jw-dev/x/chessquery/pkg/query"
)

const FileName = "test.pgn"

var queries = []query.Analyzer{
	QElo{},
	QLeastMoves{},
	QMostMoves{},
}

func main() {
	now := time.Now()
	r, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}

	runner := query.NewRunner()
	for _, q := range queries {
		runner.Add(q, "Test", query.Once)
	}

	i := 0
	for s := range pgn.Split(r) {
		game, err := pgn.Parse(s)
		if err != nil {
			log.Fatalln(err)
		}
		runner.Analyze(game)
		i += 1
	}

	fmt.Printf("Analyzed %d games with %d analyzers in %dus!\n",
		i,
		len(queries),
		time.Since(now).Microseconds())
	runner.PrintResults()
}
