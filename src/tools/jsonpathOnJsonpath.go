package tools

import (
	"strings"
)

func AreJsonPathsLinked(jsonpath1 string, jsonpath2 string) bool {

	// at this step assume that both are jsonpath

	jsonpathSanitized1 := jsonpath1
	if strings.HasSuffix(jsonpath1, ".*") {
		jsonpathSanitized1 = jsonpath1[0 : len(jsonpath1)-2]
	}

	jsonpathSanitized2 := jsonpath2
	if strings.HasSuffix(jsonpath2, ".*") {
		jsonpathSanitized2 = jsonpath2[0 : len(jsonpath2)-2]
	}

	jsonpathSlice1 := strings.Split(jsonpathSanitized1, ".")
	jsonpathSlice2 := strings.Split(jsonpathSanitized2, ".")

	var smallJsonpathSlice []string
	var longJsonpathSlice []string
	if len(jsonpathSlice1) > len(jsonpathSlice2) {
		longJsonpathSlice = jsonpathSlice1
		smallJsonpathSlice = jsonpathSlice2
	} else {
		longJsonpathSlice = jsonpathSlice2
		smallJsonpathSlice = jsonpathSlice1
	}

	for i, s := range smallJsonpathSlice {

		if s != longJsonpathSlice[i] {

			return false
		}
	}

	return true
}
