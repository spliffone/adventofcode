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

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	score := []int{0, 0}
	results := map[string][]int{
		"A X": {6 + 1, 0 + 3},
		"A Y": {3 + 1, 3 + 1},
		"A Z": {0 + 1, 6 + 2},
		"B X": {6 + 2, 0 + 1},
		"B Y": {3 + 2, 3 + 2},
		"B Z": {0 + 2, 6 + 3},
		"C X": {6 + 3, 0 + 2},
		"C Y": {3 + 3, 3 + 3},
		"C Z": {0 + 3, 6 + 1},
	}
	for sc.Scan() {
		round := sc.Text()
		points2 := results[round]
		score[0] += points2[0]
		score[1] += points2[1]

	}
	log.Println("Score player A ", strconv.Itoa(score[0]))
	log.Println("Score player B ", strconv.Itoa(score[1]))
}
