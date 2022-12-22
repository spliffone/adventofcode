package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
	knotNo    = 9
)

const (
	Up = iota
	Down
	Left
	Right
)

type Coord struct {
	Row  int
	Col  int
	Move int
}

func (c Coord) String() string {
	return fmt.Sprintf("%d#%d", c.Row, c.Col)
}

func (h Coord) Distance(o *Coord) (int, int, int, int) {
	rowDistance := 0
	colDistance := 0
	rowOffset := 0
	colOffset := 0
	if h.Row > o.Row {
		rowOffset = 1
		rowDistance = h.Row - o.Row
	} else {
		rowOffset = -1
		rowDistance = o.Row - h.Row
	}
	if h.Col > o.Col {
		colOffset = 1
		colDistance = h.Col - o.Col
	} else {
		colOffset = -1
		colDistance = o.Col - h.Col
	}
	return rowDistance, colDistance, rowOffset, colOffset
}

func (h Coord) Next(o *Coord) *Coord {
	// Check is a move necessary
	rowD, colD, rowOffset, colOffset := h.Distance(o)
	if rowD <= 1 && colD <= 1 {
		return nil
	}

	// Calc move
	if h.Row == o.Row {
		// Same horizontal line
		if h.Col < o.Col {
			// Head is left to tail -> go left
			return &Coord{Row: h.Row, Col: o.Col - 1}
		} else {
			// Head is right to tail -> go right
			return &Coord{Row: h.Row, Col: o.Col + 1}
		}
	} else if h.Col == o.Col {
		// Same vertical line
		if h.Row > o.Row {
			// Head is above
			return &Coord{Row: o.Row + 1, Col: o.Col}
		} else {
			// Head is below
			return &Coord{Row: o.Row - 1, Col: o.Col}
		}
	} else {
		// Diagonal
		return &Coord{Row: o.Row + rowOffset, Col: o.Col + colOffset}
	}
}

func (c Coord) Left() *Coord {
	return &Coord{Row: c.Row, Col: c.Col - 1, Move: Left}
}

func (c Coord) Right() *Coord {
	return &Coord{Row: c.Row, Col: c.Col + 1, Move: Right}
}

func (c Coord) Up() *Coord {
	return &Coord{Row: c.Row + 1, Col: c.Col, Move: Up}
}

func (c Coord) Down() *Coord {
	return &Coord{Row: c.Row - 1, Col: c.Col, Move: Down}
}

type Mover struct {
	knotPositions []*Coord
	visibleSteps  map[string]bool
}

func NewMover() *Mover {
	knots := make([]*Coord, 1)
	knots[0] = &Coord{
		Row: 0,
		Col: 0,
	}
	for i := 0; i < knotNo; i++ {
		knots = append(knots, &Coord{
			Row: 0,
			Col: 0,
		})
	}

	return &Mover{
		knotPositions: knots,
		visibleSteps:  map[string]bool{knots[0].String(): true},
	}
}

func (m *Mover) Move(direction string, steps string) {
	stepCount, _ := strconv.Atoi(steps)

	for s := stepCount; s > 0; s-- {
		// Move head
		head := m.knotPositions[0]
		var newHead *Coord
		switch direction {
		case "U":
			newHead = head.Up()
		case "D":
			newHead = head.Down()
		case "R":
			newHead = head.Right()
		case "L":
			newHead = head.Left()
		}
		m.knotPositions[0] = newHead

		// Move robe with all knots
		for i := 1; i < len(m.knotPositions); i++ {
			lead := m.knotPositions[i-1]
			follower := m.knotPositions[i]

			// Move tail
			newFollower := lead.Next(follower)
			if newFollower != nil {
				m.knotPositions[i] = newFollower

				// Monitor last knot position
				if i == knotNo {
					log.Println("Move tail:", newFollower.String())
					m.visibleSteps[newFollower.String()] = true
				}
			}
		}
	}
}

func (m *Mover) Result() int {
	return len(m.visibleSteps)
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
	m := NewMover()
	for row := 0; sc.Scan(); row++ {
		line := sc.Text()
		step := strings.Fields(line)
		m.Move(step[0], step[1])
	}

	log.Println("Result:", m.Result())
	if m.Result() != 6037 {
		log.Fatal("Incorrect result")
	}
}
