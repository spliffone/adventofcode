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
)

type Item struct {
	Size   int
	Name   string
	IsDir  bool
	Items  []*Item
	Parent *Item
}

func (i *Item) Add(new *Item) {
	new.Parent = i
	i.Items = append(i.Items, new)
}

func (i *Item) GetSize() int {
	if i.Size > 0 {
		return i.Size
	}

	totalSize := 0
	for _, item := range i.Items {
		totalSize += item.GetSize()
	}
	i.Size = totalSize
	return i.Size
}

func (i *Item) UpdateSize() {
	dirSize := 0
	for _, item := range i.Items {
		if !item.IsDir {
			dirSize += item.Size
		}
	}
	cur := i
	for cur != nil {
		cur.Size += dirSize
		cur = cur.Parent
	}
}

func (i *Item) GetItem(name string) *Item {
	if i == nil {
		log.Fatalln("fs not initialized")
	}
	for _, item := range i.Items {
		if item.Name == name {
			return item
		}
	}
	log.Fatalln("could not find item ", name, "in", i.Name)
	return nil
}

func ToString(i *Item, sb *strings.Builder, prefix string) {
	if i == nil {
		return
	}

	var meta string
	if i.IsDir {
		meta = fmt.Sprintf("(dir, size=%d)", i.Size)
	} else {
		meta = fmt.Sprintf("(file, size=%d)", i.Size)
	}
	sb.WriteString(fmt.Sprintf("%s%s %s\n", prefix, i.Name, meta))
	for _, c := range i.Items {
		ToString(c, sb, prefix+"  ")
	}
}

func (i *Item) String() string {
	if i == nil {
		return ""
	}
	var sb strings.Builder
	ToString(i, &sb, "")
	return sb.String()
}

type SumCandidates struct {
	Size int
}

func (c *SumCandidates) Sum100000(i *Item) {
	if i.IsDir && i.Size <= 100000 {
		log.Printf("dir %s %d", i.Name, i.Size)
		c.Size += i.Size
	}
}

func internalWalk(i *Item, f func(i *Item)) {
	if i == nil {
		return
	}
	f(i)
	for _, n := range i.Items {
		internalWalk(n, f)
	}
}

func (i *Item) Walk(f func(i *Item)) {
	internalWalk(i, f)
}

func IsCommand(i string) bool {
	return strings.HasPrefix(i, "$")
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
	mode := 0
	fs := &Item{
		Items: []*Item{{
			Name:  "/",
			IsDir: true,
		}},
	}
	current := fs
	for sc.Scan() {
		line := sc.Text()
		tokens := strings.Fields(line)
		if IsCommand(tokens[0]) {
			if mode == 1 {
				// End of ls command - calculate and propagate size
				current.UpdateSize()
			}

			if tokens[1] == "ls" {
				mode = 1
			} else if tokens[1] == "cd" {
				mode = 2
				if tokens[2] == "/" {
					current = fs.GetItem(tokens[2])
				} else if tokens[2] == ".." {
					current = current.Parent
				} else {
					current = current.GetItem(tokens[2])
				}
			}
			continue
		}
		if mode == 1 {
			if tokens[0] == "dir" {
				current.Add(&Item{
					Name:  tokens[1],
					IsDir: true,
				})
			} else {
				size, _ := strconv.Atoi(tokens[0])
				current.Add(&Item{
					Name:  tokens[1],
					Size:  size,
					IsDir: false,
				})
			}
		}
	}
	if mode == 1 {
		// End of ls command - calculate and propagate size
		current.UpdateSize()
	}

	// Walk dirs
	log.Println(fs.String())
	var c SumCandidates
	fs.Walk(c.Sum100000)
	log.Println("Result:", c.Size)
}
