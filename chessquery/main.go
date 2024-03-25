package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jw-dev/x/chessquery/pkg/pgn"
)

func main() {
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
