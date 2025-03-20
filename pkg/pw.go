package pkg

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	SPRING = 0
	SUMMER = 1
	AUTUMN = 2
	WINTER = 3
)

const (
	LEET_BASIC      = 0
	LEET_BASIC_PLUS = 1
)

func generatePW(entry *ldap.Entry, mask string) string {
	// Regular expression pattern to match placeholders with optional modifiers
	pattern := `\{([^{}]+?)(?:#([^{}]+))?\}`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Replace placeholders with function calls
	replacedMask := re.ReplaceAllStringFunc(mask, func(match string) string {
		// Extract the placeholder and modifier (if any) without curly braces
		matches := re.FindStringSubmatch(match)
		placeholder := matches[1]
		modifier := ""
		if len(matches) > 2 {
			modifier = matches[2]
		}
		placeholderLowercase := strings.ToLower(placeholder)

		// Call the appropriate function based on the placeholder and modifier
		switch placeholderLowercase {
		case "yy", "yyyy", "m", "mm", "monthgerman", "monthenglish", "seasongerman", "seasonamerican", "seasonbritish":
			year, _ := convertDate(convertTime(entry.GetAttributeValue("pwdLastSet")), placeholder)
			return year
		default:
			// If the placeholder is not recognized, try to get it from the ldap entry
			value := entry.GetAttributeValue(placeholder)
			if modifier != "" {
				result, err := applyModifier(value, modifier)
				if err != nil {
					return value
				}
				return result
			}
			return value
		}
	})

	return replacedMask
}

func leetSpeak(input string, technique int) string {
	var leetMap map[rune]string
	switch technique {
	case LEET_BASIC:
		leetMap = getLeetBasic()
	case LEET_BASIC_PLUS:
		leetMap = getLeetBasicPlus()
	}
	// Replace characters according to the leetMap
	result := ""
	for _, char := range input {
		leetChar, found := leetMap[unicode.ToUpper(char)]
		if !found {
			leetChar, found = leetMap[char]
		}
		if found {
			// Preserve the original casing
			if unicode.IsUpper(char) {
				leetChar = strings.ToUpper(leetChar)
			} else {
				leetChar = strings.ToLower(leetChar)
			}
			result += leetChar
		} else {
			result += string(char)
		}
	}
	return result
}

// convertDate converts a date string into the desired format
func convertDate(dateString string, format string) (string, error) {
	t, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(format) {
	case "yyyy":
		return fmt.Sprintf("%d", t.Year()), nil
	case "yy":
		// Get last two numbers of the year
		year := t.Year() % 100
		return fmt.Sprintf("%02d", year), nil
	case "mm":
		return fmt.Sprintf("%02d", t.Month()), nil
	case "m":
		return fmt.Sprintf("%d", t.Month()), nil
	case "monthgerman":
		return getGermanMonths()[t.Month()-1], nil
	case "monthenglish":
		return getEnglishMonths()[t.Month()-1], nil
	case "seasongerman":
		switch t.Month() {
		case time.December, time.January, time.February:
			return getGermanSeasons()[WINTER], nil
		case time.March, time.April, time.May:
			return getGermanSeasons()[SPRING], nil
		case time.June, time.July, time.August:
			return getGermanSeasons()[SUMMER], nil
		case time.September, time.October, time.November:
			return getGermanSeasons()[AUTUMN], nil
		}
	case "seasonamerican":
		switch t.Month() {
		case time.December, time.January, time.February:
			return getAmericanSeasons()[WINTER], nil
		case time.March, time.April, time.May:
			return getAmericanSeasons()[SPRING], nil
		case time.June, time.July, time.August:
			return getAmericanSeasons()[SUMMER], nil
		case time.September, time.October, time.November:
			return getAmericanSeasons()[AUTUMN], nil
		}
	case "seasonbritish":
		switch t.Month() {
		case time.December, time.January, time.February:
			return getBritishSeasons()[WINTER], nil
		case time.March, time.April, time.May:
			return getBritishSeasons()[SPRING], nil
		case time.June, time.July, time.August:
			return getBritishSeasons()[SUMMER], nil
		case time.September, time.October, time.November:
			return getBritishSeasons()[AUTUMN], nil
		}
	default:
		return "", fmt.Errorf("unsupported format")
	}
	return "", nil
}

// Reverse returns the reversed string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// splitModifiers splits a modifier string into individual modifiers
func splitModifiers(modifierString string) []string {
	var modifiers []string
	var currentModifier strings.Builder
	parenCount := 0
	escaped := false

	for i := 0; i < len(modifierString); i++ {
		char := modifierString[i]

		// Handle escape character
		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		// Handle parentheses counting
		if !escaped {
			if char == '(' {
				parenCount++
			} else if char == ')' {
				parenCount--
			}
		}

		// Handle modifier separator
		if char == '#' && !escaped && parenCount == 0 {
			// If we have collected a modifier, add it to the list
			if currentModifier.Len() > 0 {
				modifiers = append(modifiers, currentModifier.String())
				currentModifier.Reset()
			}
			continue
		}

		// Add the character to the current modifier
		if escaped && char == '#' {
			// For escaped #, add just the # without the escape character
			currentModifier.WriteByte('#')
		} else {
			currentModifier.WriteByte(char)
		}
		escaped = false
	}

	// Add the last modifier if any
	if currentModifier.Len() > 0 {
		modifiers = append(modifiers, currentModifier.String())
	}

	return modifiers
}

func applyModifier(value, modifierString string) (string, error) {
	if modifierString == "" {
		return value, nil
	}

	// Split the modifier string into individual modifiers
	modifiers := splitModifiers(modifierString)
	result := value

	// Apply each modifier in sequence
	for _, modifier := range modifiers {
		var err error
		switch {
		case modifier == "Reverse":
			result = Reverse(result)
		case modifier == "Upper":
			result = strings.ToUpper(result)
		case modifier == "Lower":
			result = strings.ToLower(result)
		case modifier == "Title":
			result = cases.Title(language.English).String(strings.ToLower(result))
		case modifier == "Capitalize":
			if len(result) > 0 {
				result = strings.ToUpper(result[:1]) + strings.ToLower(result[1:])
			}
		case modifier == "AlternateLower":
			// Convert to alternating case starting with lowercase (e.g., "Hello" -> "hElLo")
			runes := []rune(result)
			for i := range runes {
				if i%2 == 0 {
					runes[i] = unicode.ToLower(runes[i])
				} else {
					runes[i] = unicode.ToUpper(runes[i])
				}
			}
			result = string(runes)
		case modifier == "AlternateUpper":
			// Convert to alternating case starting with uppercase (e.g., "Hello" -> "HeLlO")
			runes := []rune(result)
			for i := range runes {
				if i%2 == 0 {
					runes[i] = unicode.ToUpper(runes[i])
				} else {
					runes[i] = unicode.ToLower(runes[i])
				}
			}
			result = string(runes)
		case modifier == "LeetBasic":
			result = leetSpeak(result, LEET_BASIC)
		case modifier == "LeetBasicPlus":
			result = leetSpeak(result, LEET_BASIC_PLUS)
		case strings.HasPrefix(modifier, "Pattern(") && strings.HasSuffix(modifier, ")"):
			pattern := modifier[8 : len(modifier)-1] // Extract pattern between parentheses
			result, err = ApplyPattern(result, pattern)
			if err != nil {
				return "", fmt.Errorf("error applying pattern modifier: %v", err)
			}
		default:
			return "", fmt.Errorf("unknown modifier: %s", modifier)
		}
	}

	return result, nil
}
