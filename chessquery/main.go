package main

import (
	"fmt"
	"log"
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
	for s := range pgn.Split(r) {
		_, err := pgn.Parse(s)
		if err != nil {
			log.Fatalln(err)
		}
		i += 1
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Loaded %d games in %dus (%dms)\n",
		i, elapsed.Microseconds(), elapsed.Milliseconds())
}
