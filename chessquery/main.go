package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jw-dev/x/chessquery/pkg/pgn"
	"github.com/jw-dev/x/chessquery/pkg/query"
)

const FileName = "download/lichess-202301.pgn"

var analyzers = []query.Analyzer{
	query.Analyzer{
		Name:    "EloDiff",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			return int64(math.Abs(float64(p.Meta.WhiteElo - p.Meta.BlackElo)))
		},
	},
	query.Analyzer{
		Name:    "MaxMoves",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			return int64(len(p.Meta.Moves))
		},
	},
	query.Analyzer{
		Name:     "MinMoves",
		Cadence:  query.Once,
		Reversed: true,
		Query: func(p *query.Payload) int64 {
			return int64(len(p.Meta.Moves))
		},
	},
}

func main() {
	now := time.Now()
	r, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}

	runner := query.NewRunner(analyzers...)
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
