package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

const downloadList = "https://database.lichess.org/standard/list.txt"

// download retrieves a file from the internet, and saves it to the location
// defined by `outFile`. The file is streamed, so available memory isn't
// a concern.
func download(s string, outFile string) (n int64) {
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	resp, err := http.Get(s)
	if err != nil {
		log.Fatalln(err)
	}
	n, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return n
}

func main() {
	resp, err := http.Get(downloadList)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		source := s.Text()
		_, target := path.Split(source)
		n := download(source, target)
		fmt.Printf("Download %d bytes from %s\n", n, source)
	}
}
