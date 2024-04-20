package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type grade int

const (
	Easy     grade = 5
	Ok       grade = 4
	Hard     grade = 3
	NoRecall grade = 0
)

var rGrade = map[rune]grade{
	'e': Easy,
	'o': Ok,
	'h': Hard,
	'n': NoRecall,
}

func prompt(p string, f func(s string) bool) {
	s := ""
	for {
		fmt.Print(p)
		fmt.Scanln(&s)
		if f(s) {
			break
		}
	}
}

type card struct {
	back       string
	front      string
	easiness   float64
	repetition int
	interval   int
}

func (c *card) review(q grade) {
	if q >= Ok {
		c.repetition += 1
		m := math.Round(float64(c.interval) * c.easiness)
		c.interval = int(m)
	} else {
		c.repetition = 0
		c.interval = 1
	}
	gradeFlt := float64(q)
	c.easiness += (0.1 - (5.0-gradeFlt)*(0.08+(5.0-gradeFlt)*0.02))
	c.easiness = math.Max(c.easiness, 1.3)
}

type deck struct {
	cards []card
}

func read(p string) (d deck, err error) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for i := 1; s.Scan(); i++ {
		l := s.Text()
		c := card{easiness: 2.5}
		parts := strings.Split(l, "\t")
		if len(parts) <= 1 {
			err = fmt.Errorf("err parsing line %d of length %d", i, len(parts))
			return
		}
		c.front = parts[0]
		c.back = parts[1]
		if len(parts) == 5 {
			c.interval, _ = strconv.Atoi(parts[2])
			c.repetition, _ = strconv.Atoi(parts[3])
			c.easiness, _ = strconv.ParseFloat(parts[4], 64)
		}
		d.cards = append(d.cards, c)
	}
	return
}

func (d deck) write(p string) (err error) {
	f, err := os.Create(p)
	if err != nil {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, c := range d.cards {
		f := "%s\t%s\t%d\t%d\t%f\n"
		l := fmt.Sprintf(f, c.front, c.back, c.repetition, c.interval, c.easiness)
		w.WriteString(l)
	}
	w.Flush()
	return
}

func (d deck) study() {
	l := len(d.cards)
	for i := range d.cards {
		card := &d.cards[i]
		fmt.Printf("(%d of %d) %s ", i+1, l, card.front)
		fmt.Scanln()
		fmt.Println(card.back)

		s := ""
		ok := false
		for !ok {
			fmt.Print("(e)asy (o)k (h)ard (n)o-recall (q)uit? ")
			fmt.Scanln(&s)
			s = strings.ToLower(s)
			switch {
			case s == "q":
				return
			case s == "o" || s == "h" || s == "e" || s == "n":
				r, _ := utf8.DecodeRuneInString(s)
				grade := rGrade[r]
				card.review(grade)
				ok = true
			}
		}
	}
}

func main() {
	s := "test.deck"
	deck, err := read(s)
	if err != nil {
		log.Fatal(err)
	}

	deck.study()

	if err = deck.write(s); err != nil {
		log.Fatal(err)
	}
}
