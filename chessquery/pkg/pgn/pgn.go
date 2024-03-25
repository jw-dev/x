// Package pgn implements parsing of PGN files.
package pgn

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

type meta struct {
	key   string
	value string
}

func (m meta) asInt() int {
	s := strings.Trim(m.value, " \"")
	v, _ := strconv.Atoi(s)
	return v
}

type Result struct {
	Event    string
	Site     string
	White    string
	Black    string
	WhiteElo int
	BlackElo int
	Moves    []string
}

func extractMeta(s string) (m meta, ok bool) {
	if len(s) < 3 {
		return
	}
	if s[0] != '[' || s[len(s)-1] != ']' {
		return
	}
	p := strings.Split(s[1:len(s)-1], " ")
	m.key = p[0]
	m.value = strings.Join(p[1:], " ")
	ok = true
	return
}

func extractMoves(s string) (r []string) {
	if s[:2] != "1." {
		return
	}
	t := strings.Split(s, " ")
	for _, v := range t {
		if !unicode.IsDigit(rune(v[0])) {
			r = append(r, v)
		}
	}
	return
}

func Parse(s string) (r Result, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		meta, ok := extractMeta(line)
		if ok {
			switch meta.key {
			case "Event":
				r.Event = meta.value
			case "White":
				r.White = meta.value
			case "Black":
				r.Black = meta.value
			case "WhiteElo":
				r.WhiteElo = meta.asInt()
			case "BlackElo":
				r.BlackElo = meta.asInt()
			}
			continue
		}
		if unicode.IsDigit(rune(line[0])) {
			r.Moves = extractMoves(line)
		}
	}
	return
}
