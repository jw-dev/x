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

func main() {
	now := time.Now()
	r, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}

	runner := query.NewRunner()
	runner.Add("EloDiff", QEloDiff, query.Once)
	runner.Add("MaxNumMoves", QMostMoves, query.Once)
	runner.Add("MinNumMoves", QLeastMoves, query.Once)

	i := 0
	for s := range pgn.Split(r) {
		game, err := pgn.Parse(s)
		if err != nil {
			log.Fatalln(err)
		}
		runner.Analyze(game)
		i += 1
	}

	res := runner.Results()
	fmt.Printf("Analyzed %d games with %d analyzers in %dus\n", i, len(res), time.Since(now).Microseconds())
	for _, r := range res {
		fmt.Printf("%s... %d (%s)\n", r.Name, r.Score, r.Link)
	}
}
