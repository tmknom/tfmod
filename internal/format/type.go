package format

import (
	"strings"
)

const (
	TextFormat = "text"
	JsonFormat = "json"
)

func SupportType() string {
	types := []string{TextFormat, JsonFormat}
	return strings.Join(types, ", ")
}
