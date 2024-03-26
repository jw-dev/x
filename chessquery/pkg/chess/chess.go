package chess

import (
	"errors"
)

var (
	ErrRange = errors.New("value given was out of bounds")
)

var defaultBoard Position = Position{
	whiteMove: true,
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

// Move represents a move on a Chess board.
type Move struct {
	FromCol int
	ToCol   int
	FromRow int
	ToRow   int
	Promote PieceType
}

const None = -1

func newMove() *Move {
	return &Move{
		FromCol: None,
		ToCol:   None,
		FromRow: None,
		ToRow:   None,
	}
}

type Position struct {
	whiteMove bool
	pieces    [64]Piece
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

func isMoveRune(r rune) bool {
	switch {
	case r == 'K':
		return true
	case r == 'Q':
		return true
	case r == 'R':
		return true
	case r == 'B':
		return true
	case r == 'N':
		return true
	case r == 'x':
		return true
	case r >= '1' && r <= '9':
		return true
	case r >= 'a' && r <= 'h':
		return true
	default:
		return false
	}
}

var (
	ErrInvalid    = errors.New("invalid algebraic notation")
	ErrImpossible = errors.New("given move is impossible with given board")
)

// Parse parses a move in algebraic notation (e.g., "e4" or "Qxa8" etc.) and returns a move
// for the corresponding position. Returns an error if the move is not in standard notation
// or is ambiguous.
func (p *Position) Parse(a string) (*Move, error) {
	if len(a) == 0 {
		return nil, ErrInvalid
	}
	m := newMove()
	// Castling
	if a == "O-O" || a == "O-O-O" {
		m.FromCol = 4
		m.FromRow = 0
		m.ToCol = 6
		if !p.whiteMove {
			m.FromRow = 7
		}
		if a == "O-O-O" {
			m.ToCol = 2
		}
		m.ToRow = m.FromRow
		return m, nil
	}
	// How long is the move string? e.g. gxh8=Q# <- we want to get out 'h8' here
	end := 0
	for end < len(a) && isMoveRune(rune(a[end])) {
		end += 1
	}
	if end < 2 {
		return nil, ErrInvalid
	}
	m.ToRow = int(a[end-1] - '1')
	m.ToCol = int(a[end-2] - 'a')
	// Determine the source piece, whether or not we are capturing
	// and analyse any disambiguation info we have been given
	source := Pawn
	capturing := false
	for i := 0; i < len(a) && i < end-2; i++ {
		r := a[i]
		switch {
		case r == 'K':
			source = King
		case r == 'Q':
			source = Queen
		case r == 'N':
			source = Knight
		case r == 'B':
			source = Bishop
		case r == 'R':
			source = Rook
		case r == 'x':
			capturing = true
		case r >= '1' && r <= '9':
			m.FromRow = int(r - '1')
		case r >= 'a' && r <= 'h':
			m.FromCol = int(r - 'a')
		}
	}
	if m.FromCol != None && m.ToCol != None && m.FromRow != None && m.ToRow != None {
		return m, nil
	}

	// Piece disambiguation
	switch source {
	case King:
		// For king, we just need to check the surrounding tiles to find the king
		// This can never be ambiguous
		for _, r := range []int{-1, 0, 1} {
			for _, c := range []int{-1, 0, 1} {
				if r == 0 && c == 0 {
					continue
				}
				piece, _ := p.At(m.ToCol+c, m.ToRow+r)
				if piece.Type == King {
					m.FromCol = m.ToCol + c
					m.FromRow = m.ToRow + r
					return m, nil
				}
			}
		}

		return nil, ErrImpossible
	}
	_ = source
	_ = capturing
	return m, nil
}
