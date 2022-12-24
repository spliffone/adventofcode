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

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	op := 0
	screen := make([]int, 6*40)
	screen[0] = 1

	for sc.Scan() && op < 239 {
		fields := strings.Fields(sc.Text())
		if fields[0] == "noop" {
			op += 1
			screen[op] = screen[op-1]
		} else if fields[0] == "addx" {
			v, _ := strconv.Atoi(fields[1])
			op += 1
			screen[op] = screen[op-1]
			op += 1
			screen[op] = screen[op-1] + v
		}
	}
	cycle := 0
	var sb strings.Builder
	for row := 0; row < 6; row++ {

		for col := 0; col < 40; col++ {
			sprite := screen[cycle]
			if col >= sprite-1 && col <= sprite+1 {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}

			cycle++

		}
		log.Println(sb.String())
		sb.Reset()
	}
}
