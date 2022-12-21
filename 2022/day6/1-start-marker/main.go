package main

import (
	"bufio"
	"log"
	"os"
)

const (
	inputFile = "input.txt"
)

func isStartSequence(a byte, b byte, c byte, d byte) bool {
	return a != b && a != c && a != d && b != c && b != d && c != d
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	// Read initial stacks description
	result := -1
	for sc.Scan() {
		line := sc.Text()

		for i := 0; i < len(line)-3; i++ {
			if isStartSequence(line[i], line[i+1], line[i+2], line[i+3]) {
				result = i + 4
				break
			}
		}

	}
	log.Println("Result:", result)
}
