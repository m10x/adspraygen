package pkg

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-ldap/ldap/v3"
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
				switch strings.ToLower(modifier) {
				case "reverse":
					return reverseString(value)
				case "leetbasic":
					return leetSpeak(value, LEET_BASIC)
				case "leetbasicplus":
					return leetSpeak(value, LEET_BASIC_PLUS)
				}
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

func reverseString(input string) string {
	// Convert string to a slice of runes
	runes := []rune(strings.ToLower(input))
	length := len(runes)
	// Swap characters from beginning and end
	for i := 0; i < length/2; i++ {
		runes[i], runes[length-1-i] = runes[length-1-i], runes[i]
	}
	// Convert back to string and capitalize the first letter
	reversedString := string(runes)
	if length > 0 {
		reversedString = string(unicode.ToUpper(runes[0])) + string(runes[1:])
	}
	return reversedString
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
