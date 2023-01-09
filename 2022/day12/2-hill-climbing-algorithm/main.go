package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFile = "input.txt"
)

type Pos struct {
	X int
	Y int
}

type Node struct {
	Name         string
	Value        rune
	Position     Pos
	Visited      bool
	Predecessor  *Node
	Layer        int
	adjacentList []*Node
}

type Graph struct {
	Nodes map[string]*Node
	Grid  [][]*Node
}

type Queue struct {
	A []*Node
}

func GetKey(x int, y int) string {
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

func CanMove(s *Node, n *Node) bool {
	return n.Value >= s.Value || n.Value == s.Value-1
}

func (g *Graph) GetAdjacents(n *Node) []*Node {
	var result []*Node
	// Up
	if n.Position.Y > 0 {
		neighbor := g.Grid[n.Position.Y-1][n.Position.X]
		if CanMove(n, neighbor) {
			result = append(result, neighbor)
		}
	}
	// Down
	if n.Position.Y < len(g.Grid)-1 {
		neighbor := g.Grid[n.Position.Y+1][n.Position.X]
		if CanMove(n, neighbor) {
			result = append(result, neighbor)
		}
	}
	// Left
	if n.Position.X > 0 {
		neighbor := g.Grid[n.Position.Y][n.Position.X-1]
		if CanMove(n, neighbor) {
			result = append(result, neighbor)
		}
	}
	// Right
	if n.Position.X < len(g.Grid[n.Position.Y])-1 {
		neighbor := g.Grid[n.Position.Y][n.Position.X+1]
		if CanMove(n, neighbor) {
			result = append(result, neighbor)
		}
	}
	return result
}

func (g *Graph) Print() {
	for r := 0; r < len(g.Grid); r++ {
		var sb strings.Builder
		for c := 0; c < len(g.Grid[r]); c++ {
			node := g.Grid[r][c]

			if node.Visited {
				sb.WriteString(fmt.Sprintf("*%c ", node.Value))
			} else {
				sb.WriteString(fmt.Sprintf(" %c ", node.Value))
			}
		}
		fmt.Println(sb.String())
	}
}

func NewNode(x int, y int, v rune) *Node {
	name := GetKey(x, y)
	return &Node{
		Name:    name,
		Value:   v,
		Visited: false,
		Layer:   -1,
		Position: Pos{
			X: x,
			Y: y,
		},
		Predecessor:  nil,
		adjacentList: make([]*Node, 0),
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("%s: %c", n.Name, n.Value)
}

func (q *Queue) Enqueue(n *Node) {
	q.A = append(q.A, n)
}

func (q *Queue) Dequeue() *Node {
	if len(q.A) == 0 {
		log.Fatal("Queue is empty")
	}
	n := q.A[0]
	q.A = q.A[1:]
	return n
}

func (q *Queue) IsEmpty() bool {
	return len(q.A) == 0
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	graph := Graph{Nodes: make(map[string]*Node)}
	var startNode, endNode *Node

	colLen := 0
	for row := 0; sc.Scan(); row++ {
		line := sc.Text()
		if line == "" {
			continue
		}
		graph.Grid = append(graph.Grid, make([]*Node, 0))

		for col, c := range line {
			if col > colLen {
				colLen = col
			}

			// Create node
			node := NewNode(col, row, c)
			graph.Nodes[node.Name] = node
			graph.Grid[row] = append(graph.Grid[row], node)

			// Remember start/end node
			if c == 'S' {
				node.Value = 'a'
				log.Println("Start: row=", row, "col=", col)
			} else if c == 'E' {
				startNode = node
				startNode.Value = 'z'
				log.Println("End: row=", row, "col=", col)
			}
		}

		log.Println("Read row", row)
	}
	log.Println("----------------------------")

	// Search end point
	var queue Queue
	queue.Enqueue(startNode)
	startNode.Visited = true

	for !queue.IsEmpty() && endNode == nil {
		cur := queue.Dequeue()

		// Check
		if cur.Value == 'a' {
			log.Println("Found")
			endNode = cur
			break
		}

		for _, n := range graph.GetAdjacents(cur) {
			if !n.Visited {
				n.Visited = true
				n.Predecessor = cur
				queue.Enqueue(n)
			}
		}
	}

	graph.Print()
	// Count steps
	cur := endNode

	steps := 0
	for cur != nil {
		log.Printf("Step %d: %s", steps, cur.String())
		cur = cur.Predecessor
		steps++
	}
	log.Println("Result steps: ", steps-1)
}
