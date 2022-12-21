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
	score2 := []int{0, 0}
	results := map[string][]int{
		"A X": {3 + 1, 3 + 1},
		"A Y": {0 + 1, 6 + 2},
		"A Z": {6 + 1, 0 + 3},
		"B X": {6 + 2, 0 + 1},
		"B Y": {3 + 2, 3 + 2},
		"B Z": {0 + 2, 6 + 3},
		"C X": {0 + 3, 6 + 1},
		"C Y": {6 + 3, 0 + 2},
		"C Z": {3 + 3, 3 + 3},
	}
	for sc.Scan() {
		round := sc.Text()
		points2 := results[round]
		score2[0] += points2[0]
		score2[1] += points2[1]

		actionA := int(round[0]) - 64
		actionB := int(round[2]) - 87
		points1 := []int{0, 0}

		score[0] += actionA
		score[1] += actionB
		if actionA == actionB {
			// Draw
			score[0] += 3
			score[1] += 3
			points1[0] = actionA + 3
			points1[1] = actionB + 3
		} else if actionA-1 == actionB || actionA+2 == actionB {
			// Win for A
			score[0] += 6
			score[1] += 0
			points1[0] = actionA + 6
			points1[1] = actionB + 0
		} else {
			// Win for B
			score[0] += 0
			score[1] += 6
			points1[0] = actionA + 0
			points1[1] = actionB + 6
		}
		if points1[0] != points2[0] || points1[1] != points2[1] {
			log.Fatalf("Problem the points should be equal %v <> %v", points1, points2)
		}
	}
	log.Println("Score player A ", strconv.Itoa(score[0]))
	log.Println("Score player B ", strconv.Itoa(score[1]))

	log.Println("Score player A ", strconv.Itoa(score2[0]))
	log.Println("Score player B ", strconv.Itoa(score2[1]))
}
