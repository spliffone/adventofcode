package main

import (
	"bufio"
	"log"
	"os"
	"sort"
)

const (
	inputFile = "input.txt"
)

func calcIndex(c byte) int {
	v := int(c)
	if v < 97 {
		// Uppercase range 65..90
		return v - 38
	} else {
		// Uppercase range 97..122
		return v - 96
	}
}

func find(a []byte, b []byte, c []byte, x int, y int, z int) byte {
	if a[x] == b[y] && b[y] == c[z] {
		return a[x]
	}
	if a[x] < b[y] || a[x] < c[z] {
		return find(a, b, c, x+1, y, z)
	}
	if b[y] < a[x] || b[y] < c[z] {
		return find(a, b, c, x, y+1, z)
	}
	if c[z] < a[x] || c[z] < b[y] {
		return find(a, b, c, x, y, z+1)
	}
	return 0
}

func main() {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	index, sum := 0, 0
	group := []string{"", "", ""}
	for sc.Scan() {
		group[index] = sc.Text()
		index++
		if index == 3 {
			index = 0
			m1 := []byte(group[0])
			sort.Slice(m1, func(i, j int) bool {
				return m1[i] < m1[j]
			})
			m2 := []byte(group[1])
			sort.Slice(m2, func(i, j int) bool {
				return m2[i] < m2[j]
			})
			m3 := []byte(group[2])
			sort.Slice(m3, func(i, j int) bool {
				return m3[i] < m3[j]
			})
			b := find(m1, m2, m3, 0, 0, 0)
			log.Printf("%c ", b)
			sum += calcIndex(b)
		}

	}
	log.Println("Sum: ", sum)
}
