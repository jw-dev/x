package chess

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

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

var (
	ErrInvalid    = errors.New("invalid algebraic notation")
	ErrImpossible = errors.New("given move is impossible with given board")
)

var defaultBoard *Position = FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

var pieceChar = map[rune]PieceType{
	'K': King,
	'Q': Queen,
	'N': Knight,
	'R': Rook,
	'B': Bishop,
	'P': Pawn,
}

/*
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
*/
func FEN(fen string) *Position {
	p := &Position{whiteMove: true}
	r := 0
	c := 0
	for i := 0; i < len(fen) && fen[i] != ' '; i++ {
		cr := rune(fen[i])
		if cr == '/' {
			r += 1
			c = 0
			continue
		}
		if cr >= '0' && cr <= '8' {
			c += int(cr - '0')
			continue
		}
		color := Black
		if unicode.IsUpper(cr) {
			color = White
		}
		cr = unicode.ToUpper(cr)
		if typ, ok := pieceChar[cr]; ok {
			p.setPiece(c, 7-r, MakePiece(color, typ))
		}
		c += 1
	}
	return p
}

func isMoveRune(r rune) bool {
	if _, ok := pieceChar[r]; ok {
		return true
	}
	return r == 'x' || (r >= '1' && r <= '9') || (r >= 'a' && r <= 'h')
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

func (p Position) String() string {
	b := strings.Builder{}
	for r := 7; r >= 0; r-- {
		b.WriteString(fmt.Sprintf("%c ", '1'+r))
		for c := 0; c < 8; c++ {
			piece := p.At(c, r)
			for key, value := range pieceChar {
				if value == piece.Type {
					c := key
					if piece.Color == Black {
						c = unicode.ToLower(c)
					}
					b.WriteString(fmt.Sprintf("%c ", c))
					break
				}
			}
			if piece.Type == Empty {
				b.WriteString(". ")
			}
		}
		b.WriteRune('\n')
	}
	b.WriteString("  ")
	for c := 0; c < 8; c++ {
		b.WriteString(fmt.Sprintf("%c ", 'A'+c))
	}
	return b.String()
}

func Default() Position { return *defaultBoard }

func (pos *Position) At(x, y int) (p Piece) {
	return pos.pieces[y*8+x]
}

func (ps *Position) setPiece(x, y int, p Piece) {
	ps.pieces[y*8+x] = p
}

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
	source := Empty
	capturing := false
	for i := 0; i < len(a) && i < end-2; i++ {
		r := a[i]
		src, ok := pieceChar[rune(r)]
		switch {
		case ok:
			source = src
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
				piece := p.At(m.ToCol+c, m.ToRow+r)
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
			piece := p.At(c, m.ToRow)
			if piece.Type == Queen && piece.Color == color {
				m.FromCol = c
				m.FromRow = m.ToRow
				return m, nil
			}
		}
		// Vertical.
		for r := 0; r < 8; r++ {
			piece := p.At(m.ToCol, r)
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
					piece := p.At(c, r)
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
