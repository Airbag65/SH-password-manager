package artistics

import (
	"io"
	"os"
	"strings"
)


func LoadTitle() string {
	var builder strings.Builder

	file, err := os.Open("artistics/title.txt")
	if err != nil {
		return ""
	}
	
	content, err := io.ReadAll(file) 
	if err != nil {
		return ""
	}

	builder.WriteString(string(content))

	return builder.String()
}
