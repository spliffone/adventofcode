package main

import (
	"io/ioutil"
	"log"
	"sort"
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
	Key    string
}

func NewItem(input string) *Item {
	return &Item{
		Items:  SplitList(input),
		IsList: true,
		Key:    input,
	}
}

func NewItemMixed(input string) *Item {
	return &Item{
		Items:  []string{input},
		IsList: false,
	}
}

func Compare(a *Item, b *Item) int {
	aCopy := make([]string, len(a.Items))
	copy(aCopy, a.Items)
	bCopy := make([]string, len(b.Items))
	copy(bCopy, b.Items)

	for len(aCopy) > 0 && len(bCopy) > 0 {
		// Pop left
		left := aCopy[0]
		aCopy = aCopy[1:]
		// Pop right
		right := bCopy[0]
		bCopy = bCopy[1:]

		aNumber, aErr := strconv.Atoi(left)
		bNumber, bErr := strconv.Atoi(right)

		if aErr == nil && bErr == nil {
			// Both are numbers
			if aNumber < bNumber {
				return success
			} else if aNumber > bNumber {
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

	if len(aCopy) < len(bCopy) {
		return success
	} else if len(aCopy) > len(bCopy) {
		return failed
	} else {
		return 0
	}
}

func Reduce(packages []*Item) []*Item {
	known := make(map[string]*Item, 0)
	result := make([]*Item, 0)
	for _, p := range packages {
		k := strings.Join(p.Items, ",")
		if _, exists := known[k]; !exists {
			known[k] = p
			result = append(result, p)
		}
	}
	return result
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	// Append additional divider packets
	lines = append(lines, []string{"[[2]]", "[[6]]"}...)

	var packages []*Item
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		packages = append(packages, NewItem(lines[i]))
	}

	// Sort items
	sort.Slice(packages, func(i, j int) bool {
		return Compare(packages[i], packages[j]) < 0
	})
	// Reduce
	packages = Reduce(packages)
	// Reverse
	Reverse(packages)

	log.Println("Packages:", len(packages))
	multiplies := []int{0, 0}
	for i, item := range packages {
		log.Printf("%s", item.Key)
		if item.Key == "[[2]]" {
			multiplies[0] = i + 1
		} else if item.Key == "[[6]]" {
			multiplies[1] = i + 1
		}
	}

	log.Println("Result Matching lines", multiplies[0]*multiplies[1])
}
