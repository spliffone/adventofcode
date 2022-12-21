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

func isStartSequence(s string) bool {
	log.Println("Sequence:", s, "Size:", len(s))
	var hits uint32
	hits = 0
	for i, c := range s {
		bit := int(c) - 97
		log.Printf("%d: %c bit-mask %s", i, c, strconv.FormatInt(1<<bit, 2))
		if hits&(1<<bit) == 0 {
			hits |= 1 << bit
		} else {
			return false
		}
	}

	return true
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

		for i := 0; i < len(line)-14; i++ {
			if isStartSequence(line[i : i+14]) {
				result = i + 14
				break
			}
		}

	}
	log.Println("Result:", result)
}
