package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

type Machine struct {
	relevantCycles []int
	cycle          int
	value          []*CycleValue
}

type CycleValue struct {
	Cycle       int
	Value       int
	Instruction string
}

func NewMachine() *Machine {
	values := make([]*CycleValue, 1)
	values[0] = &CycleValue{
		Value: 1,
		Cycle: 0,
	}
	relevant := []int{20, 60, 100, 140, 180, 220}
	sort.Sort(sort.Reverse(sort.IntSlice(relevant)))
	return &Machine{
		relevantCycles: relevant,
		cycle:          0,
		value:          values,
	}
}

func (m *Machine) Result() int {
	sum := 0
	cycleIndex := 0
	for i := len(m.value) - 1; i >= 0 && cycleIndex < len(m.relevantCycles); i-- {
		findCycle := m.relevantCycles[cycleIndex]
		if m.value[i].Cycle < findCycle {
			log.Printf("Cycle %d Signal: %d", findCycle, m.value[i].Value*findCycle)
			sum += m.value[i].Value * findCycle
			// Search next
			cycleIndex++
		}
	}

	return sum
}

func (m *Machine) ProcessCommand(command string) {
	fields := strings.Fields(command)

	if fields[0] == "noop" {
		m.cycle++
	} else {
		v, _ := strconv.Atoi(fields[1])
		m.cycle += 2
		m.value = append(m.value, &CycleValue{
			Value:       m.value[len(m.value)-1].Value + v,
			Cycle:       m.cycle,
			Instruction: command,
		})
	}
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	// Read grid
	m := NewMachine()
	for row := 0; sc.Scan(); row++ {
		m.ProcessCommand(sc.Text())
	}

	log.Println("Result:", m.Result())
}
