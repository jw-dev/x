package main

import (
	"fmt"
	"math"
	"os"
	"strings"
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
	{
		Name:    "MostEdge",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			i := int64(0)
			for _, s := range p.Meta.Moves {
				move, err := p.CurrentPosition.Parse(s)
				if move == nil {
					fmt.Println(s)
					panic(err)
				}
				if move.FromCol == 0 || move.FromCol == 7 {
					i += 1
				}
			}
			return i
		},
	},
	{
		Name:    "MostKingMoves",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			i := int64(0)
			for _, s := range p.Meta.Moves {
				if s[0] == 'K' || s[0] == 'k' {
					i += 1
				}
			}
			return i
		},
	},
	{

		Name:    "MostQueenMoves",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			i := int64(0)
			for _, s := range p.Meta.Moves {
				if s[0] == 'Q' || s[0] == 'q' {
					i += 1
				}
			}
			return i
		},
	},
	{
		Name:    "MostCastleChecks",
		Cadence: query.Once,
		Query: func(p *query.Payload) int64 {
			i := int64(0)
			for _, s := range p.Meta.Moves {
				if (s[len(s)-1] == '+' || s[len(s)-1] == '#') && strings.HasPrefix(s, "O-O") {
					i += 1
				}
			}
			return i
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
		fmt.Printf("%-20s%-10d%s\n", r.Name, r.Score, r.Link)
	}
}
