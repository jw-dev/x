package chess

import (
	"errors"
	"unicode"
)

var (
	ErrRange = errors.New("value given was out of bounds")
)

var defaultBoard *Position = FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

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

/*
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
*/
func FEN(fen string) *Position {
	p := &Position{whiteMove: true}
	r := 0
	c := 0
	for i := 0; i < len(fen) && fen[i] != ' '; i++ {
		cr := fen[i]
		if cr == '/' {
			r += 1
			c = 0
			continue
		}
		if cr >= '0' && cr <= '8' {
			c += int(cr - '0')
			continue
		}
		piece := Piece{}
		if unicode.IsUpper(rune(cr)) {
			piece.Color = White
		}
		switch unicode.ToUpper(rune(cr)) {
		case 'K':
			piece.Type = King
		case 'Q':
			piece.Type = Queen
		case 'R':
			piece.Type = Rook
		case 'N':
			piece.Type = Knight
		case 'B':
			piece.Type = Bishop
		case 'P':
			piece.Type = Pawn
		}
		if piece.Type != Empty {
			p.setPiece(c, 7-r, piece)
		}
		c += 1
	}
	return p
}

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

func Default() Position { return *defaultBoard }

func (pos *Position) At(x, y int) (p Piece, err error) {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		err = ErrRange
		return
	}
	p = pos.pieces[y*8+x]
	return
}

func (ps *Position) setPiece(x, y int, p Piece) {
	ps.pieces[y*8+x] = p
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
	color := White
	if !p.whiteMove {
		color = Black
	}
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
				if piece.Type == King && piece.Color == color {
					m.FromCol = m.ToCol + c
					m.FromRow = m.ToRow + r
					return m, nil
				}
			}
		}
		return nil, ErrImpossible
	case Queen:
		// Check all of the possible vectors for a queen.
		// Horizontal.
		for c := 0; c < 8; c++ {
			piece, _ := p.At(c, m.ToRow)
			if piece.Type == Queen && piece.Color == color {
				m.FromCol = c
				m.FromRow = m.ToRow
				return m, nil
			}
		}
		// Vertical.
		for r := 0; r < 8; r++ {
			piece, _ := p.At(m.ToCol, r)
			if piece.Type == Queen && piece.Color == color {
				m.FromCol = m.ToCol
				m.FromRow = r
				return m, nil
			}
		}
		// Diagonal.
		for _, dr := range []int{-1, 1} {
			for _, dc := range []int{-1, 1} {
				r := m.ToRow + dr
				c := m.ToCol + dc
				for {
					if r < 0 || c < 0 || r >= 8 || c >= 8 {
						break
					}
					piece, _ := p.At(c, r)
					if piece.Type == Queen && piece.Color == color {
						m.FromCol = c
						m.FromRow = r
						return m, nil
					}
					r += dr
					c += dc
				}
			}
		}
	}
	_ = source
	_ = capturing
	return m, nil
}
