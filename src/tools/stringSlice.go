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

func StrSliceXor(references []string, strings []string) []string {
	var notIn []string
	for _, str := range strings {
		exist := false
		for _, s := range references {
			if str == s {
				exist = true
				break
			}
		}
		if !exist {
			notIn = append(notIn, str)
		}
	}
	return notIn
}
