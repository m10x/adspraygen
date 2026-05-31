package cmd

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/m10x/adspraygen/pkg"

	"github.com/spf13/cobra"
)

var (
	patternValue string
	patternsFile string
	wordsFile    string
	numbersFile  string
	specialsFile string
	patternOut   string
	patternLimit int
)

var (
	numberTokens     = []string{"{YY}", "{YYYY}", "1", "2", "3", "12", "123"}
	specialTokens    = []string{"!", ".", "#", "-", "_"}
	wordTokens       = []string{"{monthsBritish}", "{seasonsBritish}", "{sn}", "{givenName}"}
	placeholderRegex = regexp.MustCompile(`(\[WORD\]|\[NUMBER\]|\[SPECIAL\]|\[EMPTY\])`)
)

type patternToken struct {
	kind  string
	value string
}

var patternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "Generate password masks from generic placeholder templates",
	Long:  "Implements the same behavior as generatePatterns.py with [WORD], [NUMBER], and [SPECIAL] placeholders.",
	Run: func(cmd *cobra.Command, args []string) {
		if (patternValue == "" && patternsFile == "") || (patternValue != "" && patternsFile != "") {
			pkg.PrintFatal("Specify exactly one of --pattern or --patterns-file")
		}

		// Load or use default words
		words := wordTokens
		if wordsFile != "" {
			custom := readNonEmptyLines(wordsFile)
			if len(custom) == 0 {
				pkg.PrintFatal("No words loaded. Check --words-file")
			}
			words = custom
		}

		// Load or use default numbers
		numbers := numberTokens
		if numbersFile != "" {
			custom := readNonEmptyLines(numbersFile)
			if len(custom) == 0 {
				pkg.PrintFatal("No numbers loaded. Check --numbers-file")
			}
			numbers = custom
		}

		// Load or use default specials
		specials := specialTokens
		if specialsFile != "" {
			custom := readNonEmptyLines(specialsFile)
			if len(custom) == 0 {
				pkg.PrintFatal("No specials loaded. Check --specials-file")
			}
			specials = custom
		}

		patterns := []string{patternValue}
		if patternsFile != "" {
			patterns = readNonEmptyLines(patternsFile)
			if len(patterns) == 0 {
				pkg.PrintFatal("No patterns loaded. Check --patterns-file")
			}
		}

		var writer *bufio.Writer
		var outFile *os.File
		if patternOut != "" {
			f, err := os.Create(patternOut)
			if err != nil {
				pkg.PrintFatal(err.Error())
			}
			outFile = f
			writer = bufio.NewWriter(f)
			defer func() {
				_ = writer.Flush()
				_ = outFile.Close()
			}()
		}

		produced := 0
		for _, patt := range patterns {
			tokens := expandPattern(patt)
			produceCombinations(tokens, words, numbers, specials, func(value string) bool {
				if patternLimit > 0 && produced >= patternLimit {
					return false
				}
				if writer != nil {
					_, err := writer.WriteString(value + "\n")
					if err != nil {
						pkg.PrintFatal(err.Error())
					}
				} else {
					println(value)
				}
				produced++
				return true
			})
			if patternLimit > 0 && produced >= patternLimit {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(patternCmd)

	patternCmd.Flags().StringVar(&patternValue, "pattern", "", "Single pattern string")
	patternCmd.Flags().StringVar(&patternsFile, "patterns-file", "", "File with one pattern per line")
	patternCmd.Flags().StringVar(&wordsFile, "words-file", "", "File with one word per line (uses defaults if not provided)")
	patternCmd.Flags().StringVar(&patternOut, "out", "", "Optional output file")
	patternCmd.Flags().IntVar(&patternLimit, "limit", 0, "Maximum number of generated outputs (0 = unlimited)")
	patternCmd.Flags().StringVar(&numbersFile, "numbers-file", "", "File with one number per line (uses defaults if not provided)")
	patternCmd.Flags().StringVar(&specialsFile, "specials-file", "", "File with one special character per line (uses defaults if not provided)")

	patternCmd.MarkFlagsMutuallyExclusive("pattern", "patterns-file")
	patternCmd.MarkFlagsOneRequired("pattern", "patterns-file")
}

func readNonEmptyLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		pkg.PrintFatal(err.Error())
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		pkg.PrintFatal(err.Error())
	}

	return lines
}

func buildwords(word []string) []string {
	combined := append([]string{}, word...)
	combined = append(combined, wordTokens...)

	seen := make(map[string]bool)
	out := make([]string, 0, len(combined))
	for _, token := range combined {
		if !seen[token] {
			seen[token] = true
			out = append(out, token)
		}
	}
	return out
}

func expandPattern(pattern string) []patternToken {
	// Find all placeholder occurrences ([WORD], [NUMBER], [SPECIAL], [EMPTY]) in the input pattern.
	matches := placeholderRegex.FindAllStringIndex(pattern, -1)
	// If no placeholders exist, the whole string is plain text.
	if len(matches) == 0 {
		return []patternToken{{kind: "TEXT", value: pattern}}
	}

	// Build an ordered token stream by alternating literal text and placeholder tokens.
	tokens := make([]patternToken, 0, len(matches)*2+1)
	last := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		// Add text between the previous match and the current placeholder.
		if start > last {
			tokens = append(tokens, patternToken{kind: "TEXT", value: pattern[last:start]})
		}
		// Add the placeholder token itself.
		tokens = append(tokens, patternToken{kind: "PH", value: pattern[start:end]})
		last = end
	}
	// Add trailing text after the final placeholder, if present.
	if last < len(pattern) {
		tokens = append(tokens, patternToken{kind: "TEXT", value: pattern[last:]})
	}

	return tokens
}

func produceCombinations(tokens []patternToken, words []string, numbers []string, specials []string, emit func(string) bool) {
	valuesByIndex := make([][]string, len(tokens))
	typesByIndex := make([]string, len(tokens))

	for i, tok := range tokens {
		switch {
		case tok.kind == "TEXT":
			valuesByIndex[i] = []string{tok.value}
		case tok.value == "[WORD]":
			valuesByIndex[i] = words
			typesByIndex[i] = "WORD"
		case tok.value == "[NUMBER]":
			valuesByIndex[i] = numbers
			typesByIndex[i] = "NUMBER"
		case tok.value == "[SPECIAL]":
			valuesByIndex[i] = specials
			typesByIndex[i] = "SPECIAL"
		case tok.value == "[EMPTY]":
			valuesByIndex[i] = []string{""}
		default:
			valuesByIndex[i] = []string{tok.value}
		}
	}

	seen := map[string]map[string]bool{
		"WORD":    {},
		"NUMBER":  {},
		"SPECIAL": {},
	}
	chosen := make([]string, len(tokens))

	var walk func(int) bool
	walk = func(idx int) bool {
		if idx == len(tokens) {
			return emit(strings.Join(chosen, ""))
		}

		typ := typesByIndex[idx]
		for _, candidate := range valuesByIndex[idx] {
			if typ != "" {
				if seen[typ][candidate] {
					continue
				}
				seen[typ][candidate] = true
			}

			chosen[idx] = candidate
			if !walk(idx + 1) {
				if typ != "" {
					delete(seen[typ], candidate)
				}
				return false
			}
			if typ != "" {
				delete(seen[typ], candidate)
			}
		}
		return true
	}

	walk(0)
}
