package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

type Stack struct {
	stack []byte
}

func BuildStack() *Stack {
	return &Stack{
		stack: make([]byte, 0),
	}
}

func (s *Stack) Push(c byte) {
	s.stack = append(s.stack, c)
}

func (s *Stack) PushRange(c []byte) {
	s.stack = append(s.stack, c...)
}

func (s *Stack) Pop() byte {
	size := len(s.stack)
	if size <= 0 {
		log.Fatalln("invalid operation - stack is empty")
	}
	last := s.stack[size-1]
	s.stack = s.stack[:size-1]
	return last
}

func (s *Stack) PopRange(r int) []byte {
	size := len(s.stack)
	if size <= 0 {
		log.Fatalln("invalid operation - stack is empty")
	}
	last := s.stack[size-r:]
	s.stack = s.stack[:size-r]
	return last
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
	var init []string
	var stackNumbers string
	for sc.Scan() {
		line := sc.Text()
		init = append(init, line)
		if strings.TrimSpace(line) == "" {
			stackNumbers = strings.TrimSpace(init[len(init)-2])
			// Remove last empty line and stack numbers
			init = init[:len(init)-2]
			break
		}
	}
	if stackNumbers == "" {
		log.Fatalln("unable to parse input stack description")
	}

	// Parse stack description
	var stacks []*Stack
	for _, _ = range strings.Fields(stackNumbers) {
		stacks = append(stacks, BuildStack())
	}

	// Init stacks
	for lineNo := len(init) - 1; lineNo >= 0; lineNo-- {
		stackNo := 0
		line := init[lineNo]

		for i := 0; i < len(line); i += 4 {
			if line[i] == '[' {
				stacks[stackNo].Push(line[i+1])
			}
			stackNo++
		}
	}
	log.Printf("Stacks %v", stacks)

	// Read moves
	var re = regexp.MustCompile(`^move (?P<amount>\d+) from (?P<source>\d+) to (?P<dest>\d+)`)
	for sc.Scan() {
		line := sc.Text()
		matches := re.FindStringSubmatch(line)
		amount, _ := strconv.Atoi(matches[re.SubexpIndex("amount")])
		source, _ := strconv.Atoi(matches[re.SubexpIndex("source")])
		source--
		dest, _ := strconv.Atoi(matches[re.SubexpIndex("dest")])
		dest--
		c := stacks[source].PopRange(amount)
		stacks[dest].PushRange(c)
	}
	log.Printf("Stacks %v", stacks)
	result := ""
	for i, s := range stacks {
		top := s.Pop()
		log.Printf("Top element stack %d %c", i+1, top)
		result = fmt.Sprintf("%s%c", result, top)
	}
	log.Println("Result:", result)
}
