package util

import (
	"fmt"
	"strings"
)


func RightPad(original string, length int) string {
	paddingNeeded := length - len(original)

	paddedString := fmt.Sprintf("%s%s", original, strings.Repeat(" ", paddingNeeded))

	return paddedString
}

func LeftPad(original string, length int) string {
	paddingNeeded := length - len(original)

	paddedString := fmt.Sprintf("%s%s", strings.Repeat(" ", paddingNeeded), original)

	return paddedString
}
