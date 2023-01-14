package tui

import (
	"encoding/json"
	"strings"
)

func getShortMethodName(methodName string) string {
	return methodName[strings.LastIndex(methodName, ".")+1:]
}

func countLines(s string) int {
	return strings.Count(s, "\n")
}

func checkJSON(s string) error {
	return json.Unmarshal([]byte(s), &struct{}{})
}
