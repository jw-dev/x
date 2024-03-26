package chess

import (
	"errors"
)

var (
	ErrRange = errors.New("value given was out of bounds")
)

var defaultBoard Position = Position{
	pieces: [64]Piece{
		0: {1, 4}, 1: {1, 1}, 6: {0, 1}, 7: {0, 4}, 8: {1, 3}, 9: {1, 1}, 14: {0, 1}, 15: {0, 3}, 16: {1, 2}, 17: {1, 1}, 22: {0, 1}, 23: {0, 2}, 24: {1, 5}, 25: {1, 1}, 30: {0, 1}, 31: {0, 5}, 32: {1, 6}, 33: {1, 1}, 38: {0, 1}, 39: {0, 6}, 40: {1, 2}, 41: {1, 1}, 46: {0, 1}, 47: {0, 2}, 48: {1, 3}, 49: {1, 1}, 54: {0, 1}, 55: {0, 3}, 56: {1, 4}, 57: {1, 1}, 62: {0, 1}, 63: {0, 4},
	},
}

const (
	Black Color = iota
	White
)

const (
	Empty PieceType = iota
	Pawn
	Bishop
	Knight
	Rook
	Queen
	King
)

type Color uint8
type PieceType uint8

type Piece struct {
	Color Color
	Type  PieceType
}

func MakePiece(c Color, t PieceType) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}

type Position struct {
	pieces [64]Piece
}

func Default() Position { return defaultBoard }

func (pos *Position) At(x, y int) (p Piece, err error) {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		err = ErrRange
		return
	}
	p = pos.pieces[x*8+y]
	return
}
