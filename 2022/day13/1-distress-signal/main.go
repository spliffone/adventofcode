package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
	success   = 1
	failed    = -1
)

func SplitList(input string) []string {
	var result []string
	var open []int

	for index := 0; index < len(input); index++ {
		if input[index] == '[' {
			open = append(open, index)
		} else if input[index] == ']' {
			// Finished - reached end of list
			if len(open) == 1 {
				listStart := open[0] + 1
				// Ignore empty list
				if input[listStart] != ']' {
					result = append(result, input[listStart:index])
				}
				break
			}

			// Remove last element
			open = open[0 : len(open)-1]

		} else if input[index] == ',' {
			// Only handle elements in the same scope
			if len(open) == 1 {
				listStart := open[0] + 1
				result = append(result, input[listStart:index])
				open[0] = index
			}
		}
	}
	return result
}

type Item struct {
	Items  []string
	IsList bool
}

func NewItem(input string) *Item {
	return &Item{
		Items:  SplitList(input),
		IsList: true,
	}
}

func NewItemMixed(input string) *Item {
	return &Item{
		Items:  []string{input},
		IsList: false,
	}
}

func Compare(a *Item, b *Item) int {
	log.Printf("Compare %v vs %v", a.Items, b.Items)

	for len(a.Items) > 0 && len(b.Items) > 0 {
		// Pop left
		left := a.Items[0]
		a.Items = a.Items[1:]
		// Pop right
		right := b.Items[0]
		b.Items = b.Items[1:]

		aNumber, aErr := strconv.Atoi(left)
		bNumber, bErr := strconv.Atoi(right)

		if aErr == nil && bErr == nil {
			log.Printf("Compare %d vs %d", aNumber, bNumber)
			// Both are numbers
			if aNumber < bNumber {
				return success
			} else if aNumber > bNumber {
				log.Println("Right side is smaller, so inputs are not in the right order")
				return failed
			}
		} else if aErr != nil && bErr != nil {
			// Both are lists
			result := Compare(NewItem(left), NewItem(right))
			if result != 0 {
				return result
			}
		} else if aErr == nil && bErr != nil {
			// Left is number right is list
			result := Compare(NewItemMixed(left), NewItem(right))
			if result != 0 {
				return result
			}
		} else if aErr != nil && bErr == nil {
			// Left is list, right is number
			result := Compare(NewItem(left), NewItemMixed(right))
			if result != 0 {
				return result
			}
		}
	}

	if len(a.Items) < len(b.Items) {
		return success
	} else if len(a.Items) > len(b.Items) {
		return failed
	} else {
		return 0
	}
}

func main() {
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	matching := 0
	pair := 1
	var list []int
	for i := 0; i < len(lines); i += 3 {
		log.Println("== Compare pair", pair)

		result := Compare(NewItem(lines[i]), NewItem(lines[i+1]))
		if result == 1 {
			log.Println("true")
			list = append(list, pair)
			matching += pair
		} else {
			log.Println("false")
		}

		pair++
	}
	log.Printf("%v", list)
	log.Println("Result Matching lines", matching)
}
