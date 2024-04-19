package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
algorithm SM-2 is
    input:  user grade q
            repetition number n
            easiness factor EF
            interval I
    output: updated values of n, EF, and I

    if q ≥ 3 (correct response) then
        if n = 0 then
            I ← 1
        else if n = 1 then
            I ← 6
        else
            I ← round(I × EF)
        end if
        increment n
    else (incorrect response)
        n ← 0
        I ← 1
    end if

    EF ← EF + (0.1 − (5 − q) × (0.08 + (5 − q) × 0.02))
    if EF < 1.3 then
        EF ← 1.3
    end if

    return (n, EF, I)
*/

type deck struct {
	cards []card
}

type card struct {
	back     string
	front    string
	grade    int
	rep      int
	interval int
}

func read(p string) (d deck, err error) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			continue
		}
		c := card{}
		parts := strings.Split(l, "\t")
		len := len(parts)
		switch len {
		case 4:
			if v, err := strconv.Atoi(parts[2]); err != nil {
				c.rep = v
			}
			if v, err := strconv.Atoi(parts[3]); err != nil {
				c.interval = v
			}
			fallthrough
		case 2:
			c.front = parts[0]
			c.back = parts[1]
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
		l := fmt.Sprintf("%s\t%s\t%d\t%d\n", c.front, c.back, c.rep, c.interval)
		w.WriteString(l)
	}
	w.Flush()
	return
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

func (d deck) study() {
	l := len(d.cards)
	for i, card := range d.cards {
		fmt.Printf("(%d of %d) %s", i, l, card.front)
		fmt.Scanln()
		fmt.Println(card.back)
		prompt("(y)es (h)ard (n)o (q)uit? ", func(s string) bool {
			return true
		})
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
