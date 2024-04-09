package chess_test

import (
	"testing"

	"github.com/jw-dev/x/chessquery/pkg/chess"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	pos := chess.Default()
	exp := chess.MakePiece(chess.White, chess.Rook)
	p := pos.At(0, 0)
	require.Equal(t, exp, p, "piece did not match")
}

func TestParseCastling(t *testing.T) {
	s := "O-O"
	exp := &chess.Move{
		FromCol: 4,
		ToCol:   6,
	}
	pos := chess.Default()
	act, err := pos.Parse(s)
	require.Nil(t, err)
	require.Equal(t, exp, act)
}

func TestParseLongCastle(t *testing.T) {
	s := "O-O-O"
	exp := &chess.Move{
		FromCol: 4,
		ToCol:   2,
	}
	pos := chess.Default()
	act, err := pos.Parse(s)
	require.Nil(t, err)
	require.Equal(t, exp, act)
}

func TestParseUnambiguous(t *testing.T) {
	s := "Qc4xe6#"
	exp := &chess.Move{
		FromCol: 2,
		ToCol:   4,
		FromRow: 3,
		ToRow:   5,
	}
	pos := chess.Default()
	act, err := pos.Parse(s)
	require.Nil(t, err)
	require.Equal(t, exp, act)
}

func TestParsePromotion(t *testing.T) {
	s := "cxd8=Q#" // c pawn takes d on rank 8, promotes to Queen, gives mate
	exp := &chess.Move{
		FromCol: 2,
		ToCol:   3,
		// FromRow: 6, FIXME
		ToRow: 7,
		// Promote: chess.Queen, FIXME
	}

	pos := chess.Default()
	act, err := pos.Parse(s)
	require.Nil(t, err)
	require.Equal(t, exp, act)
}

func TestParseKingMove(t *testing.T) {
	s := "Kf2"
	exp := &chess.Move{
		FromCol: 4,
		ToCol:   5,
		FromRow: 0,
		ToRow:   1,
	}
	pos := chess.Default()
	act, err := pos.Parse(s)
	require.Nil(t, err)
	require.Equal(t, exp, act)
}

func BenchmarkXxx(b *testing.B) {
	s := "cxd8=Q#" // c pawn takes d on rank 8, promotes to Queen, gives mate
	pos := chess.Default()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pos.Parse(s)
	}
}
