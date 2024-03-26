package chess_test

import (
	"testing"

	"github.com/jw-dev/x/chessquery/pkg/chess"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	pos := chess.Default()
	exp := chess.MakePiece(chess.White, chess.Rook)
	p, err := pos.At(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exp, p, "piece did not match")
}
