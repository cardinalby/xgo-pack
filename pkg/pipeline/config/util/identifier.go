package util

import (
	"slices"
	"strings"
)

func reverseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	slices.Reverse(parts)
	return strings.Join(parts, ".")
}

func ModulePathToIdentifier(modulePath string) string {
	pathParts := strings.Split(modulePath, "/")
	if len(pathParts) >= 1 {
		pathParts[0] = reverseDomain(pathParts[0])
	}
	return strings.Join(pathParts, ".")
}

func IdentifierLastPart(identifier string) string {
	parts := strings.Split(identifier, ".")
	return parts[len(parts)-1]
}

func IdentifierWithoutLastPart(identifier string) string {
	parts := strings.Split(identifier, ".")
	if len(parts) == 1 {
		return identifier
	}
	return strings.Join(parts[:len(parts)-1], ".")
}
