package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	NoColor = 0
	Red     = 1
	Yellow  = 2
	Green   = 3
	Cyan    = 4
	Blue    = 5
	Magenta = 6
)

// formatTimestamp returns a formatted timestamp for logging
func formatTimestamp() string {
	return time.Now().Format("15:04:05")
}

func Print(msg string, c int) {
	timestamp := color.HiBlackString("[%s]", formatTimestamp())
	var prefix string

	switch c {
	case Red:
		prefix = color.RedString("ERROR")
	case Yellow:
		prefix = color.YellowString("WARN ")
	case Green:
		prefix = color.GreenString("OK   ")
	case Cyan:
		prefix = color.CyanString("INFO ")
	case Blue:
		prefix = color.BlueString("DEBUG")
	case Magenta:
		prefix = color.MagentaString("TRACE")
	default:
		prefix = "     "
	}

	// Format the message with proper indentation for multiline
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if i == 0 {
			fmt.Printf("%s %s %s\n", timestamp, prefix, line)
		} else if line != "" {
			fmt.Printf("%s %s %s\n", timestamp, strings.Repeat(" ", 5), line)
		}
	}
}

func PrintFatal(msg string) {
	Print(msg, Red)
	if !strings.HasSuffix(msg, "\n") {
		fmt.Println()
	}
	os.Exit(1)
}

// PrintSuccess prints a success message in green with a checkmark
func PrintSuccess(msg string) {
	Print("✓ "+msg, Green)
}

// PrintError prints an error message in red with an X
func PrintError(msg string) {
	Print("✗ "+msg, Red)
}

// PrintWarning prints a warning message in yellow with a warning symbol
func PrintWarning(msg string) {
	Print("⚠ "+msg, Yellow)
}

// PrintInfo prints an info message in cyan with an info symbol
func PrintInfo(msg string) {
	Print("ℹ "+msg, Cyan)
}

// PrintDebug prints a debug message in blue
func PrintDebug(msg string) {
	Print("⚙ "+msg, Blue)
}

// ReadMaskFile reads a file and returns all non-empty, non-comment lines as masks.
func ReadMaskFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open mask file: %w", err)
	}
	defer f.Close()

	var masks []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		masks = append(masks, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading mask file: %w", err)
	}
	return masks, nil
}
