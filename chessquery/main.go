package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jw-dev/x/chessquery/pkg/pgn"
)

func main() {
	r, err := os.Open("test.pgn")
	if err != nil {
		panic(err)
	}

	i := 0
	start := time.Now()
	for range pgn.Split(r) {
		i += 1
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Loaded %d games in %dus (%dms)\n",
		i, elapsed.Microseconds(), elapsed.Milliseconds())
}
