package query

import (
	"fmt"
	"math"

	"github.com/jw-dev/x/chessquery/pkg/chess"
	"github.com/jw-dev/x/chessquery/pkg/pgn"
)

// Cadence is a type to identify when a query should be run.
// A cadence of `Once` is called once, at the beginning of a game.
// `EveryPly` is called every single ply.
// `AtEnd` is called once, after the end of the game.
type Cadence int8

const (
	Once Cadence = iota
	EveryPly
	AtEnd
)

// Payload is the data an Analyzer receives from the current chess
// game that is running.
type Payload struct {
	Meta            pgn.Result
	LastMove        chess.Move
	CurrentPosition chess.Position
}

func newPayload(m pgn.Result) Payload {
	return Payload{
		Meta: m,
	}
}

// Analyzer takes in a Payload and returns an int64 denoting
// the `score` of the current position. Higher score = higher
// weighting for the position.
type Analyzer interface {
	Analyze(p *Payload) int64
}

type analyzerMeta struct {
	name  string
	cad   Cadence
	score int64
	site  string
	Analyzer
}

type Runner struct {
	analyzers []analyzerMeta
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Add(a Analyzer, name string, cad Cadence) {
	r.analyzers = append(r.analyzers, analyzerMeta{
		name:     name,
		cad:      cad,
		score:    math.MinInt64,
		Analyzer: a,
	})
}

func (r *Runner) Analyze(g pgn.Result) {
	// Only `Meta` type analyzers are supported right now...
	for i := range r.analyzers {
		a := &r.analyzers[i]
		payload := newPayload(g)
		score := a.Analyze(&payload)
		if score > a.score {
			a.score = score
			a.site = g.Site
		}
	}
}

// TEMP FUNCTION
func (r *Runner) PrintResults() {
	for _, a := range r.analyzers {
		fmt.Printf("%s: %d (%s)\n", a.name, a.score, a.site)
	}
}
