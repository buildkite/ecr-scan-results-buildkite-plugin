package env

import (
	"os"
	"strings"
)

// ParseWithPrefix will read env vars with a matching prefix into a slice.
func ParseWithPrefix(prefix string) []string {
	result := []string{}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		val := parts[1]

		if len(strings.TrimSpace(val)) > 0 && strings.HasPrefix(key, prefix) {
			result = append(result, val)
		}
	}
	return result
}
