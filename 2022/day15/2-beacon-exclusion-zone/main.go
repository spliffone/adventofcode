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
					// Remove gaps
					inRange = inRange || r[1]+1 == r2[0]

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

func GetOpenSpot(scans []Scan, rowMax int) (int, int) {

	// Find ranges
	for r := 0; r <= rowMax; r++ {
		var ranges [][2]int
		for _, scan := range scans {
			distance := scan.ManhattanDistanceToBeacon()
			xRange := distance - scan.RowDistance(r)
			if xRange >= 0 {
				s := Max(0, scan.SensorX-xRange)
				e := Min(rowMax, scan.SensorX+xRange)
				ranges = append(ranges, [2]int{s, e})
			}
		}
		ranges = MergeRanges(ranges)
		//fmt.Printf("%d: %v\n", r, ranges)
		if len(ranges) > 1 {
			return ranges[0][1] + 1, r
		}
	}
	return 0, 0
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

	// x, y - 0 and no larger than 4000000
	// tuning frequency, which can be found by multiplying its x coordinate by 4000000 and then adding its y coordinate
	// frequency := x * 4000000 + y
	x, y := GetOpenSpot(scans, 4000000)
	frequency := x*4000000 + y
	log.Println("frequency", frequency)
}
