package aspect

import (
	"strings"
)

type options []string

func (o options) Has(option string) bool {
	for _, opt := range o {
		if opt == option {
			return true
		}
	}
	return false
}

func parseTag(tag string) (string, options) {
	parts := strings.Split(tag, ",")
	return parts[0], options(parts[1:])
}
