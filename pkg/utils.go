package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	NoColor = 0
	Red     = 1
	Yellow  = 2
	Green   = 3
	Cyan    = 4
)

func Print(msg string, c int) {
	if c == Red {
		msg = color.RedString("[ERR] ") + msg
	} else if c == Yellow {
		msg = color.YellowString("[!] ") + msg
	} else if c == Green {
		msg = color.GreenString("[+] ") + msg
	} else if c == Cyan {
		msg = color.CyanString("[*] ") + msg
	}

	fmt.Print(msg)
}

func PrintFatal(msg string) {
	Print(msg, Red)
	if !strings.HasSuffix(msg, "\n") {
		fmt.Println()
	}
	os.Exit(1)
}
