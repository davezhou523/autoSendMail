package helper

import "strings"

func StrSplit(str string, length int) []string {
	str = strings.Join(strings.Fields(str), "")
	strs := []rune(str)
	c := len(strs)
	arr := []string{}
	if length < 1 || length >= c {
		return []string{str}
	}
	for i := 0; i < c; i += length {
		arr = append(arr, string(strs[i:i+length]))
	}
	return arr
}

func MbStrpos(haystack, needle string) int {
	index := strings.Index(haystack, needle)
	if index == -1 || index == 0 {
		return index
	}
	pos := 0
	total := 0
	reader := strings.NewReader(haystack)
	for {
		_, size, err := reader.ReadRune()
		if err != nil {
			return -1
		}
		total += size
		pos++
		// got it
		if total == index {
			return pos
		}
	}
}
