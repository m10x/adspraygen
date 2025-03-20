package pkg

import (
	"fmt"
	"strings"
)

// PatternReplacement represents a single pattern replacement rule
type PatternReplacement struct {
	From string
	To   string
}

// unescapePattern handles escape sequences in pattern strings
func unescapePattern(s string) string {
	replacements := map[string]string{
		`\;`: ";",
		`\(`: "(",
		`\>`: ">",
		`\)`: ")",
		`\\`: "\\",
	}

	result := s
	for escaped, unescaped := range replacements {
		result = strings.ReplaceAll(result, escaped, unescaped)
	}
	return result
}

// ParsePattern parses a pattern string in the format "from>to;from2>to2"
func ParsePattern(pattern string) ([]PatternReplacement, error) {
	var replacements []PatternReplacement
	var currentPart strings.Builder
	escaped := false

	// Split pattern into individual replacements handling escaped characters
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]

		if escaped {
			currentPart.WriteByte(char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == ';' {
			// Process the current part
			if err := processPart(&replacements, currentPart.String()); err != nil {
				return nil, err
			}
			currentPart.Reset()
			continue
		}

		currentPart.WriteByte(char)
	}

	// Process the last part
	if currentPart.Len() > 0 {
		if err := processPart(&replacements, currentPart.String()); err != nil {
			return nil, err
		}
	}

	return replacements, nil
}

// processPart processes a single pattern part and adds it to the replacements slice
func processPart(replacements *[]PatternReplacement, part string) error {
	// Split each part into from and to
	fromTo := strings.SplitN(part, ">", 2)
	if len(fromTo) != 2 {
		return fmt.Errorf("invalid pattern format in '%s', expected 'from>to'", part)
	}

	from := strings.TrimSpace(unescapePattern(fromTo[0]))
	to := strings.TrimSpace(unescapePattern(fromTo[1]))
	if from == "" {
		return fmt.Errorf("'from' pattern cannot be empty in '%s'", part)
	}

	*replacements = append(*replacements, PatternReplacement{
		From: from,
		To:   to,
	})

	return nil
}

// ApplyPattern applies the pattern replacements to the input string
func ApplyPattern(input string, pattern string) (string, error) {
	replacements, err := ParsePattern(pattern)
	if err != nil {
		return "", err
	}

	result := input
	for _, replacement := range replacements {
		result = strings.ReplaceAll(result, replacement.From, replacement.To)
	}

	return result, nil
}
