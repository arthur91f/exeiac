package tools

import "regexp"

type ReplaceOperation struct {
	Regex     *regexp.Regexp
	ReplaceBy string
}

func CreateReplaceOperation(regex string, replace string) (
	o ReplaceOperation, err error) {
	o.Regex, err = regexp.Compile(regex)
	if err != nil {
		return
	}
	o.ReplaceBy = replace
	return
}

func (r ReplaceOperation) Replace(i string) (o string) {
	return r.Regex.ReplaceAllString(i, r.ReplaceBy)
}
