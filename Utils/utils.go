package Utils

import "strings"
import "regexp"

func UpperConcat(s ...string) string {
	result := strings.ToLower(s[0])
	for i := 1; i < len(s); i++ {
		result += strings.Title(s[i])
	}
	return result
}

func RemoveUnderscore(s string) string {
	re := regexp.MustCompile(`_+`)
	replaced := re.ReplaceAllString(s, "")
	return replaced
}

func ReSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && i < len(match) && name != "" {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}
