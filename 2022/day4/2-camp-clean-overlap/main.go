package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

func extractRange(i string) []int {
	numbers := strings.Split(i, "-")
	s, _ := strconv.Atoi(numbers[0])
	e, _ := strconv.Atoi(numbers[1])
	return []int{s, e}
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	fullSubsets := 0
	for sc.Scan() {
		pair := sc.Text()
		sections := strings.Split(pair, ",")
		aRange := extractRange(sections[0])
		bRange := extractRange(sections[1])
		// Is subset:
		// First range is subset ( (rs1 >= rs2 && rs1 <= re2) || (re1 >= rs2 && re1 <= re2) ) ||
		// Second range is subset ( (rs2 >= rs1 && rs2 <= re1) || (re2 >= rs1 && re2 <= re1) )
		if (aRange[0] >= bRange[0] && aRange[0] <= bRange[1]) || (aRange[1] >= bRange[0] && aRange[1] <= bRange[1]) ||
			(bRange[0] >= aRange[0] && bRange[0] <= aRange[1]) || (bRange[1] >= aRange[0] && bRange[1] <= aRange[1]) {
			fullSubsets++
			log.Println("Subset ", pair)
		}

	}
	log.Println("Full subsets ", fullSubsets)
}
