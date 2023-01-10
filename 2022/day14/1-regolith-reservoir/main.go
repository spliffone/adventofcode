package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

func Min(c, o int) int {
	if c < o {
		return c
	}
	return o
}

func Max(c, o int) int {
	if c > o {
		return c
	}
	return o
}

type Grid struct {
	grid [][]int
	xMin int
}

func NewGrid(xMin, xMax, yMax int) *Grid {
	rows := yMax
	cols := xMax - xMin
	grid := make([][]int, rows+1)
	for row := 0; row < len(grid); row++ {
		grid[row] = make([]int, cols+1)
	}
	return &Grid{
		grid: grid,
		xMin: xMin,
	}
}

func (g *Grid) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(g.xMin, "..", g.xMin+len(g.grid[0])))
	for i, rows := range g.grid {
		sb.WriteString(fmt.Sprintf("%4d %v", i, rows))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *Grid) RelativeX(x int) int {
	return x - g.xMin
}

func (g *Grid) CalcOffset(x, y, targetX, targetY int) []int {
	o := []int{0, 0}
	if y == targetY {
		if targetX < x {
			o[0] = -1
		} else {
			o[0] = 1
		}

	} else if x == targetX {
		if targetY < y {
			o[1] = -1
		} else {
			o[1] = 1
		}
	}
	return o
}

func (g *Grid) AddStones(p [][]int) {
	for _, path := range p {
		// Start
		x, y := g.RelativeX(path[0]), path[1]
		g.grid[y][x] = 1
		for i, j := 2, 3; j < len(path); i, j = i+2, j+2 {
			targetX, targetY := g.RelativeX(path[i]), path[j]

			o := g.CalcOffset(x, y, targetX, targetY)
			// Set stones
			for x != targetX || y != targetY {
				y = y + o[1]
				x = x + o[0]
				g.grid[y][x] = 1
			}
		}
	}
}

const (
	Overflow = 0
)

func (g *Grid) CanMoveDown(x, y int) int {
	if y+1 < len(g.grid)-1 {
		return g.grid[y+1][x]
	} else {
		return -1
	}
}

func (g *Grid) CanMoveDiagonalLeft(x, y int) int {
	if x == 0 {
		return -1
	}

	return g.grid[y+1][x-1]
}

func (g *Grid) CanMoveDiagonalRight(x, y int) int {
	if x == len(g.grid[y+1])-1 {
		return -1
	}

	return g.grid[y+1][x+1]
}

func (g *Grid) Leak2(x, y int) int {
	startX, startY := x, y
	more := true
	for more {
		// Fall down
		v := g.CanMoveDown(x, y)
		if v < 0 {
			//
			return -1
		} else {
			// Position is free
			if v == 0 {
				y++
				continue
			}
		}
		v = g.CanMoveDiagonalLeft(x, y)
		if v < 0 {
			return -1
		} else {
			if v == 0 {
				y++
				x--
				continue
			}
		}
		v = g.CanMoveDiagonalRight(x, y)
		if v < 0 {
			return -1
		} else {
			if v == 0 {
				y++
				x++
				continue
			}
		}
		more = false
	}
	if startX != x && startY != y {
		g.grid[y][x] = 2
		return 0
	}
	return -1
}

func (g *Grid) Leak(x, y int) bool {

	more := true
	for more {
		newY := y + 1
		if newY == len(g.grid) {
			y = newY
			more = false
			break
		} else if x-1 < 0 {
			y = newY
			x--
			more = false
			break
		} else if x+1 == len(g.grid[0]) {
			y = newY
			x++
			more = false
			break
		}

		if newY < len(g.grid) && g.grid[newY][x] == 0 {
			// First: Can we move down
			y = newY
		} else if newY < len(g.grid) && x > 0 && g.grid[newY][x-1] == 0 {
			// Second: Can we move diagonal left
			y = newY
			x--
		} else if newY < len(g.grid) && x+1 < len(g.grid[0]) && g.grid[newY][x+1] == 0 {
			// Third: Can we move diagonal right
			y = newY
			x++
		} else {
			// Can't find a possible next move
			more = false
		}
	}

	if x > 0 && x < len(g.grid[0]) && y < len(g.grid) {
		g.grid[y][x] = 2
		more = true
	} else {
		more = false
		fmt.Print("Finished")
	}

	return more
}

func main() {
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	xMin, xMax := math.MaxInt32, 0
	yMax := 0
	var paths [][]int
	// Parse paths
	for i := 0; i < len(lines); i++ {
		path := lines[i]
		if path == "" {
			continue
		}
		var coords []int
		for _, pos := range strings.Split(path, " -> ") {
			positions := strings.Split(pos, ",")
			x, _ := strconv.Atoi(positions[0])
			coords = append(coords, x)
			xMin = Min(xMin, x)
			xMax = Max(xMax, x)
			y, _ := strconv.Atoi(positions[1])
			coords = append(coords, y)
			yMax = Max(yMax, y)
		}

		paths = append(paths, coords)
	}

	// Build grid
	grid := NewGrid(xMin, xMax, yMax)
	grid.AddStones(paths)
	fmt.Println(grid.String())
	var i int
	for i = 0; grid.Leak(grid.RelativeX(500), 0); i++ {
		fmt.Println("Unit:", i)
		fmt.Println(grid.String())
		if i == 22 {
			fmt.Println(grid.String())
		}
	}
	log.Println("Result Matching lines", i)
}
