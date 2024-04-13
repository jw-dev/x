// Package pgn implements parsing of PGN files.
package pgn

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func isKeyValue(s string) bool {
	return len(s) > 4 && s[0] == '[' && s[len(s)-1] == ']'
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

func (r *Result) zero() {
	r.Event = ""
	r.Site = ""
	r.White = ""
	r.Black = ""
	r.WhiteElo = 0
	r.BlackElo = 0
	r.Moves = make([]string, 0)
}

func (r *Result) extractKeyValue(s string) {
	if len(s) < 3 {
		return
	}
	if s[0] != '[' || s[len(s)-1] != ']' {
		return
	}
	mid := strings.IndexByte(s, ' ')
	if mid > -1 {
		key := s[1:mid]
		q1 := strings.IndexByte(s, '"')
		q2 := strings.LastIndexByte(s, '"')
		value := s[q1+1 : q2]
		switch key {
		case "Event":
			r.Event = value
		case "White":
			r.White = value
		case "Black":
			r.Black = value
		case "WhiteElo":
			v, err := strconv.Atoi(value)
			if err != nil {
				r.WhiteElo = v
			}
		case "BlackElo":
			v, err := strconv.Atoi(value)
			if err != nil {
				r.BlackElo = v
			}
		case "Site":
			r.Site = value
		}
	}
}

func (r *Result) extractMoves(moves string) {
	r.Moves = make([]string, 0, len(moves)/6+10)
	start := 0
	for i := 1; i < len(moves); i++ {
		c := moves[i]
		if unicode.IsLetter(rune(c)) && moves[i-1] == ' ' {
			start = i
		}
		if start > 0 && (c == ' ' || c == '?' || c == '!') {
			r.Moves = append(r.Moves, moves[start:i])
			start = 0
		}
	}
}

func ReadAll(r io.Reader, f func(r *Result)) {
	sc := bufio.NewScanner(r)
	res := new(Result)
	inPgn := true
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 {
			continue
		}
		if isKeyValue(line) {
			if !inPgn {
				f(res)
				res.zero()
				inPgn = true
			}
			res.extractKeyValue(line)
			continue
		}
		inPgn = false
		if unicode.IsNumber(rune(line[0])) {
			res.extractMoves(line)
		}
	}
}
