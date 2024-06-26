package query

import (
	"math"

	"github.com/jw-dev/x/chessquery/pkg/chess"
	"github.com/jw-dev/x/chessquery/pkg/pgn"
)

// Function is an analyzer that takes in the board state and returns
// a score. The higher the score, the more 'interesting' the position
// is, according to that analyzer.
type Query func(p *Payload) int64

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
	Meta            *pgn.Result
	LastMove        chess.Move
	CurrentPosition chess.Position
}

func newPayload(m *pgn.Result) Payload {
	return Payload{
		Meta:            m,
		CurrentPosition: chess.Default(),
	}
}

// Analyzer takes in a Payload and returns an int64 denoting
// the `score` of the current position. Higher score = higher
// weighting for the position.
type Analyzer struct {
	Name     string
	Cadence  Cadence
	Query    Query
	Reversed bool
	score    int64
	link     string
}

type Result struct {
	Name  string
	Score int64
	Link  string
}

type Runner struct {
	analyzers []Analyzer
}

func NewRunner(a ...Analyzer) *Runner {
	for i := range a {
		a[i].score = math.MinInt64
		if a[i].Reversed {
			a[i].score = math.MaxInt64
		}
	}
	return &Runner{analyzers: a}
}

func (r *Runner) Analyze(g *pgn.Result) {
	// Only `Meta` type analyzers are supported right now...
	for i := range r.analyzers {
		a := &r.analyzers[i]
		payload := newPayload(g)
		score := a.Query(&payload)
		if (!a.Reversed && score > a.score) || (a.Reversed && score < a.score) {
			a.score = score
			a.link = g.Site
		}
	}
}

func (r *Runner) Results() (res []Result) {
	for _, a := range r.analyzers {
		res = append(res, Result{
			Name:  a.Name,
			Score: a.score,
			Link:  a.link,
		})
	}
	return
}
