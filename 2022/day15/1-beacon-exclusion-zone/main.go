package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

var (
	r = regexp.MustCompile(`Sensor at x=(?P<SensorX>-?\d+), y=(?P<SensorY>-?\d+): closest beacon is at x=(?P<BeaconX>-?\d+), y=(?P<BeaconY>-?\d+).*`)
)

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type Scan struct {
	SensorX int
	SensorY int
	BeaconX int
	BeaconY int
}

func (s Scan) ManhattanDistance(x, y int) int {
	return Abs(s.SensorX-x) + Abs(s.SensorY-y)
}

func (s Scan) ManhattanDistanceToBeacon() int {
	return s.ManhattanDistance(s.BeaconX, s.BeaconY)
}

func (s Scan) String() string {
	return fmt.Sprintf("Sensor x=%d, y=%d, Beacon x=%d, y=%d, Manhatten distance=%d", s.SensorX, s.SensorY, s.BeaconX, s.BeaconY, s.ManhattanDistanceToBeacon())
}

func (s Scan) RowDistance(row int) int {
	return Abs(row - s.SensorY)
}

func MergeRanges(ranges [][2]int) [][2]int {
	foundMerge := true

Merged:
	for foundMerge {
		for i, r := range ranges {
			for j, r2 := range ranges {
				if i != j {

					inRange := r[0] >= r2[0] && r[1] <= r2[1] || r[1] >= r2[0] && r[1] <= r2[1]
					inRange = inRange || r2[0] >= r[0] && r2[1] <= r[1] || r2[1] >= r[0] && r2[1] <= r[1]

					if inRange {
						r2[0] = Min(r[0], r2[0])
						r2[1] = Max(r[1], r2[1])
						ranges[j] = r2
						// Remove
						ranges = append(ranges[:i], ranges[i+1:]...)
						continue Merged
					}
				}
			}
		}
		foundMerge = false
	}

	return ranges
}

func CalcOccupied(scans []Scan, row int) int {

	// Find ranges
	var ranges [][2]int
	beaconsCols := make(map[int]int, 0)
	for _, scan := range scans {
		distance := scan.ManhattanDistanceToBeacon()
		xRange := distance - scan.RowDistance(row)
		if xRange >= 0 {
			ranges = append(ranges, [2]int{scan.SensorX - xRange, scan.SensorX + xRange})
		}
		if scan.BeaconY == row {
			beaconsCols[scan.BeaconX] = 1
		}
	}

	fmt.Printf("Ranges before merge: %v\n", ranges)
	ranges = MergeRanges(ranges)
	fmt.Printf("Ranges after merge: %v\n", ranges)

	// Count occupied area
	count := 0
	for _, r := range ranges {
		if r[0] == r[1] {
			count++
		} else {
			count += r[1] - r[0]
            // 
			if 0 >= r[0] && 0 <= r[1] {
				count++
			}
		}

		for k := range beaconsCols {
			if r[0] <= k && k <= r[1] {
				count--
			}
		}
	}
	return count
}

func main() {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	var scans []Scan
	for i := 0; i < len(lines); i++ {
		info := lines[i]
		if info == "" {
			continue
		}
		matches := r.FindStringSubmatch(info)
		if len(matches) >= 5 {
			xSensor, _ := strconv.Atoi(matches[1])
			ySensor, _ := strconv.Atoi(matches[2])
			xBeacon, _ := strconv.Atoi(matches[3])
			yBeacon, _ := strconv.Atoi(matches[4])
			scan := Scan{
				SensorX: xSensor,
				SensorY: ySensor,
				BeaconX: xBeacon,
				BeaconY: yBeacon,
			}
			scans = append(scans, scan)
		} else {
			log.Fatalln("Unable to parse line")
		}
	}

	log.Println("how many positions cannot contain a beacon", CalcOccupied(scans, 2000000))
}
