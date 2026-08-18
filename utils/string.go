package utils

import "strings"

func NormalizeRegionName(name string) string {
	name = strings.ToUpper(strings.TrimSpace(name))

	prefixes := []string{
		"KABUPATEN ",
		"KAB. ",
		"KAB ",
		"KOTA ",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}

	return strings.TrimSpace(name)
}
