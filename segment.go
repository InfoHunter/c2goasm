package main

import (
	"fmt"
	"strings"
)

type Segment struct {
	Name       string
	Start, End int
	stack      Stack
}

func segmentEqual(a, b []Segment) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !(a[i].Name == b[i].Name && a[i].Start == b[i].Start && a[i].End == b[i].End) {
			return false
		}
	}

	return true
}

func SegmentSource(src []string) []Segment {

	segments := []Segment{}

	gatherUntilRetForSegment := -1

	for index, line := range src {

		// Find start of a subroutine
		if strings.Contains(line, "## @") {
			entryName := ExtractName(strings.Split(line, "## @")[1])

			for _, s := range segments {
				if s.Name == entryName {
					panic(fmt.Sprintf("Entry name %s already found", entryName))
				}
			}

			segments = append(segments, Segment{Name: entryName, Start: index + 1})
		}

		// Find end of a subroutine
		if strings.Contains(line, ".exit") {
			exitName := ExtractName(strings.Split(line, "## %")[1])

			isegment := -1
			var s Segment
			for isegment, s = range segments {
				if s.Name == exitName {
					break
				}
			}
			if isegment == -1 || isegment == len(segments) {
				panic(fmt.Sprintf("No entry name found for exit %s", exitName))
			}

			segments[isegment].End = index + 1 // include this line (label)

			// Gather stack information
			gatherUntilRetForSegment = isegment
		}

		if gatherUntilRetForSegment != -1 {
			if strings.Contains(line, "ret") {

				// Lines of postamble
				stackLines := src[segments[gatherUntilRetForSegment].End:index+1]

				segments[gatherUntilRetForSegment].stack = ExtractStackInfo(stackLines)

				gatherUntilRetForSegment = -1
			}
		}
	}

	return segments
}
