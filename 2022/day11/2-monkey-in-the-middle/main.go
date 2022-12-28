package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	inputFile     = "input.txt"
	StartItems    = "  Starting items: "
	Operation     = "  Operation: new = "
	Test          = "  Test: divisible by "
	SuccessAction = "    If true: throw to monkey "
	FalseAction   = "    If false: throw to monkey "
)

type Monkey struct {
	Number int
	// Starting items lists your worry level for each item the monkey
	// is currently holding in the order they will be inspected.
	items []int
	// Operation shows how your worry level changes as that monkey inspects an item.
	// (An operation like new = old * 5 means that your worry level
	// after the monkey inspected the item is five times whatever
	// your worry level was before inspection.)
	Op string
	// Test shows how the monkey uses your worry level to decide where to throw an item next.
	// If true shows what happens with an item if the Test was true.
	// If false shows what happens with an item if the Test was false.
	TestDivisor   int
	SuccessAction int
	FalseAction   int
	Inspection    int
}

func NewMonkey(number int) *Monkey {
	return &Monkey{
		Number: number,
	}
}

// When a monkey throws an item to another monkey,
// the item goes on the end of the recipient monkey's list.
func (m *Monkey) Add(item int) {
	m.items = append(m.items, item)
}

func (m *Monkey) Pop() int {
	v := m.items[0]
	m.items = m.items[1:]
	m.Inspection++
	return v
}

func (m *Monkey) SetOperation(op string) {
	m.Op = op
}

func (m *Monkey) HasItems() bool {
	return len(m.items) > 0
}

func (m *Monkey) Operation(worryLvl int) int {
	test := strings.Fields(m.Op)
	if test[0] == "old" && test[1] == "*" && test[2] == "old" {
		return worryLvl * worryLvl
	} else if test[0] == "old" {
		v, _ := strconv.Atoi(test[2])
		switch test[1] {
		case "+":
			return worryLvl + int(v)
		case "-":
			return worryLvl - int(v)
		case "*":
			return worryLvl * int(v)
		case "/":
			return worryLvl / int(v)
		}
	}
	log.Fatal("Invalid!")
	return 0
}

func (m *Monkey) Test(v int) bool {
	return true
}

func (m *Monkey) String() string {
	return fmt.Sprintf("Monkey %d (i=%d): %v", m.Number, m.Inspection, m.items)
}

func CommonModulo(monkeys []*Monkey) int {
	v := 1
	for _, m := range monkeys {
		v *= m.TestDivisor
	}
	return v
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	monkeys := make([]*Monkey, 0)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		cols := strings.Fields(line)
		if cols[0] == "Monkey" {

			number, _ := strconv.Atoi(strings.Trim(cols[1], ":"))
			m := NewMonkey(number)

			// Parse: Starting items: 79, 98
			sc.Scan()
			line = sc.Text()
			if strings.HasPrefix(line, StartItems) {
				items := strings.Fields(line[len(StartItems):])
				for _, item := range items {
					item = strings.TrimRight(item, ",")
					no, _ := strconv.Atoi(item)
					m.Add(int(no))
				}
			}
			sc.Scan()
			line = sc.Text()
			if strings.HasPrefix(line, Operation) {
				m.SetOperation(line[len(Operation):])

			}
			sc.Scan()
			line = sc.Text()
			if strings.HasPrefix(line, Test) {
				v, _ := strconv.Atoi(line[len(Test):])
				m.TestDivisor = int(v)
			}
			sc.Scan()
			line = sc.Text()
			if strings.HasPrefix(line, SuccessAction) {
				v, _ := strconv.Atoi(line[len(SuccessAction):])
				m.SuccessAction = int(v)
			}
			sc.Scan()
			line = sc.Text()
			if strings.HasPrefix(line, FalseAction) {
				v, _ := strconv.Atoi(line[len(FalseAction):])
				m.FalseAction = int(v)
			}

			monkeys = append(monkeys, m)
		}
	}

	// Multiply all test values to get a common modulo
	commonMod := CommonModulo(monkeys)

	monkeyInspections := make([]int, len(monkeys))
	for round := 1; round <= 10000; round++ {
		// The monkeys take turns inspecting and throwing items.
		// On a single monkey's turn, it inspects and throws all of
		// the items it is holding one at a time and in the order listed.
		for _, m := range monkeys {

			// Get Next Monkey and make his turn
			for m.HasItems() {
				item := m.Pop()
				newLvl := m.Operation(item)
				reduced := newLvl % commonMod

				monkeyInspections[m.Number] = monkeyInspections[m.Number] + 1

				if reduced%m.TestDivisor == 0 {
					monkeys[m.SuccessAction].Add(reduced)
				} else {
					monkeys[m.FalseAction].Add(reduced)
				}
			}
		}

		if round%1000 == 0 {
			log.Printf("Round %d", round)
			for _, m := range monkeys {
				log.Println(m.String())
			}
		}

	}

	sort.Ints(monkeyInspections)
	log.Printf("Result: %v", monkeyInspections)
	log.Printf("Result: %d", monkeyInspections[len(monkeyInspections)-1]*monkeyInspections[len(monkeyInspections)-2])

}
