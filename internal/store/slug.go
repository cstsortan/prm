package store

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)
	multiDash       = regexp.MustCompile(`-{2,}`)
)

// GenerateSlug creates a URL-safe slug from a title.
// It lowercases, normalizes unicode, replaces non-alphanumeric chars with dashes,
// and trims leading/trailing dashes.
func GenerateSlug(title string) string {
	// Replace common symbols with words before normalizing
	title = strings.ReplaceAll(title, "&", " and ")
	title = strings.ReplaceAll(title, "+", " plus ")

	// Normalize unicode and lowercase
	s := strings.ToLower(norm.NFKD.String(title))

	// Remove non-ASCII characters (accents etc.)
	var b strings.Builder
	for _, r := range s {
		if r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == ' ') {
			b.WriteRune(r)
		}
	}
	s = b.String()

	// Replace spaces and non-alphanumeric with dashes
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = multiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if s == "" {
		s = "untitled"
	}

	// Cap slug length to avoid filesystem path limits
	if len(s) > 80 {
		s = s[:80]
		s = strings.TrimRight(s, "-")
	}

	return s
}

// UniqueSlug returns a slug that doesn't collide with existing directories in parentDir.
// If "my-slug" exists, it tries "my-slug-2", "my-slug-3", etc.
// Gives up after 1000 attempts and returns the last candidate.
func UniqueSlug(parentDir, slug string) string {
	candidate := slug
	counter := 2
	for counter <= 1000 {
		if _, err := os.Stat(filepath.Join(parentDir, candidate)); os.IsNotExist(err) {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", slug, counter)
		counter++
	}
	return candidate
}
