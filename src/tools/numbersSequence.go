package tools

import (
	"fmt"
	"strconv"
	"strings"
)

type NumbersSequence string

func (s NumbersSequence) IsValid() bool {

	elements := strings.Split(string(s), ",")

	for _, element := range elements {

		if strings.Contains(element, "-") {
			rangeParts := strings.Split(element, "-")
			if len(rangeParts) == 2 {

				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err1 != nil || err2 != nil {

					return false
				}
				if start > end {

					return false
				}
			} else {

				return false
			}
		} else {

			_, err := strconv.Atoi(strings.TrimSpace(element))
			if err != nil {

				return false
			}
		}
	}

	return true
}

func (s NumbersSequence) Contains(number int) (bool, error) {
	if !s.IsValid() {

		return false, fmt.Errorf("invalid sequence number format: %s", s)
	}

	// Split the sequence by commas
	elements := strings.Split(string(s), ",")

	// Iterate over each element in the sequence
	for _, element := range elements {
		// Check if the element is a range
		if strings.Contains(element, "-") {
			rangeParts := strings.Split(element, "-")
			// Check if the start and end of the range are valid numbers
			start, _ := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			end, _ := strconv.Atoi(strings.TrimSpace(rangeParts[1]))

			// Check if the number is within the range
			if number >= start && number <= end {
				return true, nil
			}
		} else {
			// Check if the element is equal to the number
			value, err := strconv.Atoi(strings.TrimSpace(element))
			if err == nil && value == number {
				return true, nil
			}
		}
	}

	return false, nil
}
