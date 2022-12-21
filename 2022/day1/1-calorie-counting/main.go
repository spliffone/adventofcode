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

	pos, max := 0, 0
	calories := []int{0}
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			calories = append(calories, 0)
			pos++
		} else {
			cal, _ := strconv.Atoi(line)
			calories[pos] += cal
			max = Max(max, calories[pos])
		}
	}
	log.Println("Max calorie is ", max)
}

func Max(l int, r int) int {
	if l < r {
		return r
	}
	return l
}
