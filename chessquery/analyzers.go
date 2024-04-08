package main

import (
	"math"

	"github.com/jw-dev/x/chessquery/pkg/query"
)

func QEloDiff(p *query.Payload) int64 {
	return int64(math.Abs(float64(p.Meta.WhiteElo) - float64(p.Meta.BlackElo)))
}

func QMostMoves(p *query.Payload) int64 {
	return int64(len(p.Meta.Moves))
}

func QLeastMoves(p *query.Payload) int64 {
	return int64(math.MaxInt64 - len(p.Meta.Moves))
}
