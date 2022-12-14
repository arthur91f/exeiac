package tools

func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func Deduplicate(ss []string) (output []string) {
	allKeys := make(map[string]bool)
	for _, s := range ss {
		if _, ok := allKeys[s]; !ok {
			allKeys[s] = true
			output = append(output, s)
		}
	}

	return
}
