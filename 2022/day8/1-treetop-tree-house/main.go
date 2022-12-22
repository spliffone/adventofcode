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

func MoveGetMax(grid [][]int, row int, col int, v int, offset []int, cache map[int]int) int {
	row = row + offset[0]
	col = col + offset[1]
	current := grid[row][col]

	key := row + 1000*col + 1000000*offset[0] + 1000000000*offset[1]
	if max, exists := cache[key]; exists {
		return max
	}

	// End of grid
	if row == 0 || col == 0 || row == len(grid)-1 || col == len(grid[row])-1 {
		cache[key] = current
		return current
	}

	m := MoveGetMax(grid, row, col, v, offset, cache)
	// Update max
	if m > current {
		current = m
	}
	cache[key] = current
	return current
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	grid := make([][]int, 0)
	// Read grid
	for row := 0; sc.Scan(); row++ {
		line := sc.Text()
		grid = append(grid, make([]int, len(line)))
		for col := 0; col < len(line); col++ {
			v, _ := strconv.Atoi(string(line[col]))
			grid[row][col] = v
		}
	}
	// Find the visible items
	rowLen := len(grid)
	colLen := len(grid[0])

	// The outer borders are always visible
	visible := 1
	if rowLen > 1 && colLen > 1 {
		cache := make(map[int]int, 0)
		visible = rowLen*2 + (colLen*2 - 4 /* Skip corners */)
		for row := 1; row < rowLen-1; row++ {
			for col := 1; col < colLen-1; col++ {
				v := grid[row][col]
				log.Println("Check [", row, ",", col, "]")
				// Up
				if v > MoveGetMax(grid, row, col, v, []int{-1, 0}, cache) {
					visible++
					continue
				}
				// Down
				if v > MoveGetMax(grid, row, col, v, []int{1, 0}, cache) {
					visible++
					continue
				}
				// Left
				if v > MoveGetMax(grid, row, col, v, []int{0, -1}, cache) {
					visible++
					continue
				}
				// Right
				if v > MoveGetMax(grid, row, col, v, []int{0, 1}, cache) {
					visible++
					continue
				}
			}
		}
	}

	log.Println("Result:", visible)
}
