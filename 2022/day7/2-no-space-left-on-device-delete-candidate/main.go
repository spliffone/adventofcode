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

func internalWalkDir(i *Item, f func(i *Item)) {
	if i == nil {
		return
	}
	if i.IsDir {
		f(i)
	}
	for _, n := range i.Items {
		internalWalkDir(n, f)
	}
}

func (i *Item) WalkDir(f func(i *Item)) {
	internalWalkDir(i, f)
}

type DeleteCandidates struct {
	Size      int
	Candidate *Item
}

func (c *DeleteCandidates) FindBestCandidate(i *Item) {
	if 70000000 >= c.Size-i.Size+30000000 {
		log.Printf("dir %s %d", i.Name, i.Size)
		if i.Size < c.Candidate.Size {
			c.Candidate = i
		}
	}
}

func IsCommand(i string) bool {
	return strings.HasPrefix(i, "$")
}

func NewDir(name string) *Item {
	return &Item{
		Name:  name,
		IsDir: true,
	}
}

func NewFile(name string, size int) *Item {
	return &Item{
		Name:  name,
		IsDir: false,
		Size:  size,
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

	// Read initial stacks description
	mode := 0
	fs := &Item{
		Items: []*Item{NewDir("/")},
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
				current.Add(NewDir(tokens[1]))
			} else {
				size, _ := strconv.Atoi(tokens[0])
				current.Add(NewFile(tokens[1], size))
			}
		}
	}
	if mode == 1 {
		// End of ls command - calculate and propagate size
		current.UpdateSize()
	}

	// Walk dirs
	log.Println(fs.String())
	root := fs.GetItem("/")
	c := DeleteCandidates{
		Size:      root.Size,
		Candidate: root,
	}

	fs.WalkDir(c.FindBestCandidate)
	log.Println("Result:", c.Candidate.Name, "size:", c.Candidate.Size)
}
