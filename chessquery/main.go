package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/jw-dev/x/chessquery/pkg/pgn"
	"github.com/jw-dev/x/chessquery/pkg/query"
)

const FileName = "download/lichess-202301.pgn"

var analyzers = []query.Analyzer{
	{
		Name:    "EloDiff",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			return int64(math.Abs(float64(p.Meta.WhiteElo - p.Meta.BlackElo)))
		},
	},
	{
		Name:    "MaxMoves",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			return int64(len(p.Meta.Moves))
		},
	},
	{
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

	pgn.ReadAll(r, func(r *pgn.Result) {
		runner.Analyze(r)
		i += 1
	})

	res := runner.Results()
	fmt.Printf("Analyzed %d games with %d analyzers in %.2fs\n",
		i,
		len(res),
		time.Since(now).Seconds())
	for _, r := range res {
		fmt.Printf("%-15s%-10d%s\n", r.Name, r.Score, r.Link)
	}
}
