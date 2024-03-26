package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jw-dev/x/chessquery/pkg/chess"
	"github.com/jw-dev/x/chessquery/pkg/pgn"
)

func main() {
	p := chess.Default()
	fmt.Println(p)

	data, err := os.ReadFile("test.pgn")
	if err != nil {
		log.Fatalln(err)
	}
	result, err := pgn.Parse(string(data))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}
