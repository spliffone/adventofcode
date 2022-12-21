package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
)

const (
	inputFile = "input.txt"
)

type LinkedNode struct {
	Value int
	Prev  *LinkedNode
	Next  *LinkedNode
}

type LinkedList struct {
	Head *LinkedNode
	Size int
}

func BuildLinkedList() *LinkedList {
	return &LinkedList{
		Head: nil,
		Size: 0,
	}

}

func (l *LinkedList) Insert(v int) {
	if l.Head == nil {
		l.Head = &LinkedNode{
			Value: v,
		}
		l.Size = 1
		return
	}

	cur := l.Head
	for cur.Next != nil && v > cur.Value {
		cur = cur.Next
	}
	n := &LinkedNode{
		Value: v,
	}

	if v <= cur.Value {
		// Link before
		n.Prev = cur.Prev
		n.Next = cur
		if cur.Prev != nil {
			cur.Prev.Next = n
		}
		cur.Prev = n
		if cur == l.Head {
			l.Head = n
		}
	} else {
		// Link after
		n.Prev = cur
		cur.Next = n
	}
	l.Size++

	if l.Size > 3 {
		// Trim
		l.Head = l.Head.Next
		l.Head.Prev.Next = nil
		l.Head.Prev = nil
		l.Size--
	}
	log.Println(l.String())
}

func (l LinkedList) String() string {
	var b bytes.Buffer
	cur := l.Head
	for cur != nil {
		b.WriteString(strconv.Itoa(cur.Value) + " ")
		cur = cur.Next
	}
	return b.String()
}

func (l LinkedList) Sum() int {
	cur := l.Head
	sum := 0
	for cur != nil {
		sum += cur.Value
		cur = cur.Next
	}
	return sum
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	current := 0
	top := BuildLinkedList()
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			top.Insert(current)
			current = 0
		} else {
			cal, _ := strconv.Atoi(line)
			current += cal
		}
	}
	log.Println("Max calories are: ", top.String())
	log.Println("Total of top3 calories are: ", strconv.Itoa(top.Sum()))
}
