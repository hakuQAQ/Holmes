package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// 字符串数组去重
func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func SlicetoSting(silce []string) string {
	tostring := fmt.Sprintf(strings.Join(silce, ","))
	return tostring
}

func StringtoSlice(str string) []string {
	toslice := strings.Split(str, ",")
	return toslice
}

func LinestoSlice(str string) []string {
	toslice := strings.Split(str, "\n")
	return toslice
}

func StringReplce(old string, mathes []string, new string) string {
	for _, math := range mathes {
		old = strings.Replace(old, math, new, -1)
	}
	return old
}

func RegexpTitle(content string) string {
	reTitle := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	matchResults := reTitle.FindAllString(content, -1)
	var nilstring = ""
	var mathes = []string{"<title>", "</title>"}
	return StringReplce(SlicetoSting(matchResults), mathes, nilstring)
}

func RegexpBanner(content string) string {
	reTitle := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	matchResults := reTitle.FindAllString(content, -1)
	var nilstring = ""
	var mathes = []string{"<title>", "</title>"}
	return StringReplce(SlicetoSting(matchResults), mathes, nilstring)
}