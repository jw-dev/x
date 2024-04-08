package main

import (
	"math"

	"github.com/jw-dev/x/chessquery/pkg/query"
)

type QElo struct{}

func (q QElo) Analyze(p *query.Payload) int64 {
	return int64(math.Abs(float64(p.Meta.WhiteElo) - float64(p.Meta.BlackElo)))
}

type QMostMoves struct{}

func (q QMostMoves) Analyze(p *query.Payload) int64 {
	return int64(len(p.Meta.Moves))
}

type QLeastMoves struct{}

func (q QLeastMoves) Analyze(p *query.Payload) int64 {
	return int64(math.MaxInt64 - len(p.Meta.Moves))
}
