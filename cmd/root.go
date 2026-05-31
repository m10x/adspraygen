package cmd

import (
	"fmt"
	"strings"

	"github.com/m10x/adspraygen/pkg"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "v1.2.0"
	rootCmd = &cobra.Command{
		Version: version,
		Use:     "adspraygen",
		Short:   "Active Directory password spray helper toolkit",
		Long:    fmt.Sprintf("%s\nADSprayGen %s\n\nUse one of the available subcommands: gen, pattern, spray", getLogo(), version),
	}
)

func getMaskOptions() string {
	return `
Mask Placeholders
- {cn} : Full Name
- {givenName} : First Name
- {sn} : Last Name
- {sAMAccountName} : Logon Name (Pre Windows 2000)
- {rPrincipalName} : Logon Name
- {description} : Description
- {info} : Notes
- {department} : Department
- {I} : City
- {postcalCode} : Postal Code
- Last password change
    - {YYYY} : e.g. 2024
    - {YY} : e.g. 24
    - {MM} : e.g. 01
    - {M} : e.g. 1
    - {SeasonGerman} : e.g. Herbst
    - {SeasonAmerican} : e.g. Fall
    - {SeasonBritish} : e.g. Autumn
    - {MonthGerman} : e.g. Januar
    - {MonthEnglish} : e.g. January

Mask Placeholder Modifiers (append to placeholder with #, chainable with multiple #)
- #Upper            : Convert to uppercase                  e.g. {givenName#Upper} → JOHN
- #Lower            : Convert to lowercase                  e.g. {givenName#Lower} → john
- #Title            : Capitalize first letter of each word  e.g. {givenName#Title} → John Smith
- #Capitalize       : Capitalize first letter only          e.g. {givenName#Capitalize} → John
- #AlternateLower   : Alternating case, start lower         e.g. {givenName#AlternateLower} → jOhN
- #AlternateUpper   : Alternating case, start upper         e.g. {givenName#AlternateUpper} → JoHn
- #Reverse          : Reverse the string                    e.g. {givenName#Reverse} → nhoJ
- #LeetBasic        : Substitute e→3, o→0, i→1, a→4
- #LeetBasicPlus    : Substitute e→3, o→0, i→1, a→@, t→7
- #Pattern(x>y)     : Replace x with y, chain rules with ;  e.g. {sn#Pattern(e>3;a>4)} → l33tspeak`
}

func getLogo() (logo string) {
	// source: https://patorjk.com/software/taag/#p=display&v=3&f=Bloody&t=ADSprayGen
	logo = `
 ▄▄▄      ▓█████▄   ██████  ██▓███   ██▀███   ▄▄▄     ▓██   ██▓  ▄████ ▓█████  ███▄    █ 
▒████▄    ▒██▀ ██▌▒██    ▒ ▓██░  ██▒▓██ ▒ ██▒▒████▄    ▒██  ██▒ ██▒ ▀█▒▓█   ▀  ██ ▀█   █ 
▒██  ▀█▄  ░██   █▌░ ▓██▄   ▓██░ ██▓▒▓██ ░▄█ ▒▒██  ▀█▄   ▒██ ██░▒██░▄▄▄░▒███   ▓██  ▀█ ██▒
░██▄▄▄▄██ ░▓█▄   ▌  ▒   ██▒▒██▄█▓▒ ▒▒██▀▀█▄  ░██▄▄▄▄██  ░ ▐██▓░░▓█  ██▓▒▓█  ▄ ▓██▒  ▐▌██▒
 ▓█   ▓██▒░▒████▓ ▒██████▒▒▒██▒ ░  ░░██▓ ▒██▒ ▓█   ▓██▒ ░ ██▒▓░░▒▓███▀▒░▒████▒▒██░   ▓██░
 ▒▒   ▓▒█░ ▒▒▓  ▒ ▒ ▒▓▒ ▒ ░▒▓▒░ ░  ░░ ▒▓ ░▒▓░ ▒▒   ▓▒█░  ██▒▒▒  ░▒   ▒ ░░ ▒░ ░░ ▒░   ▒ ▒ 
  ▒   ▒▒ ░ ░ ▒  ▒ ░ ░▒  ░ ░░▒ ░       ░▒ ░ ▒░  ▒   ▒▒ ░▓██ ░▒░   ░   ░  ░ ░  ░░ ░░   ░ ▒░
  ░   ▒    ░ ░  ░ ░  ░  ░  ░░         ░░   ░   ░   ▒   ▒ ▒ ░░  ░ ░   ░    ░      ░   ░ ░ 
      ░  ░   ░          ░              ░           ░  ░░ ░           ░    ░  ░         ░ 
           ░                                           ░ ░                               `

	logo = strings.ReplaceAll(logo, "░", color.MagentaString("░"))
	logo = strings.ReplaceAll(logo, "▒", color.HiMagentaString("▒"))
	return
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pkg.PrintFatal(err.Error())
	}
}
