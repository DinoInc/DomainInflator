package Utils

import "strings"
import "regexp"
import "math/rand"
import "encoding/hex"
import "os"
import "path/filepath"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func UpperConcat(s ...string) string {
	result := s[0]
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

func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
