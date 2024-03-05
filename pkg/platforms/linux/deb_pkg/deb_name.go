package deb_pkg

import (
	"fmt"
	"regexp"
	"strings"
)

// must consist only of lower case letters (a-z), digits (0-9), plus (+) and minus (-) signs, and periods (.).
// They must be at least two characters long and must start with an alphanumeric character.
// https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-source

func validateDebName(debName string) error {
	regExp := regexp.MustCompile(`^[a-z0-9][a-z0-9+-.]+$`)
	if !regExp.MatchString(debName) {
		return fmt.Errorf(
			"invalid deb package name: %s, "+
				"see https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-source",
			debName,
		)
	}
	return nil
}

func toDebName(name string) string {
	name = strings.ToLower(name)
	regExp := regexp.MustCompile(`[^a-z0-9+-.]`)
	return regExp.ReplaceAllString(name, "-")
}
