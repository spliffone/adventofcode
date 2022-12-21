package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

const (
	inputFile = "input.txt"
)

type KnownChars struct {
	Known []bool
}

func BuildKnownChars() *KnownChars {
	return &KnownChars{
		// Lowercase item types a through z have priorities 1 through 26.
		// Uppercase item types A through Z have priorities 27 through 52.
		Known: make([]bool, 53),
	}
}

func (k *KnownChars) CalcIndex(c byte) int {
	v := int(c)
	if v < 97 {
		// Uppercase range 65..90
		return v - 38
	} else {
		// Uppercase range 97..122
		return v - 96
	}
}

func (k *KnownChars) Use(c byte) *KnownChars {
	i := k.CalcIndex(c)
	k.Known[i] = true
	return k
}

func (k *KnownChars) InUse(c byte) bool {
	i := k.CalcIndex(c)
	return k.Known[i]
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	priority := 0
	for sc.Scan() {
		text := sc.Text()
		mid := len(text) / 2
		compoundA := text[0:mid]
		compoundB := text[mid:]
		a := BuildKnownChars()
		b := BuildKnownChars()
		for i := 0; i < mid; i++ {
			charA := compoundA[i]
			charB := compoundB[i]
			if a.Use(charA).InUse(charB) {
				index := a.CalcIndex(charB)
				priority += index
				log.Printf("%c is in both with prio %d\n", charB, index)
				break
			}
			if b.Use(charB).InUse(charA) {
				index := a.CalcIndex(charA)
				priority += index
				log.Printf("%c is in both with prio %d\n", charA, index)
				break
			}
		}

	}
	log.Println("Priority is ", strconv.Itoa(priority))
}
