package utils

import (
	"fmt"
	"strings"
)

func ReplaceMacros(template string, macros map[string]string) string {
	for key, value := range macros {
		placeholder := fmt.Sprintf("{{%s}}", key)
		template = strings.ReplaceAll(template, placeholder, value)
	}
	return template
}
