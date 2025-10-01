package cmd

import (
	"fmt"
	"strings"

	"github.com/m10x/adspraygen/pkg"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "v1.1.1"
	rootCmd = &cobra.Command{
		Version: version,
		Use:     "adspraygen",
		Short:   getShortDescription(),
		Long:    fmt.Sprintf("%s\nADSprayGen %s\n\n%s\n%s", getLogo(), version, getShortDescription(), getMaskOptions()),
		Example: "adspraygen -d domain.local -u m10x -p m10x -s 10.10.10.10 -m 'Foobar{givenName#Reverse}{MonthGerman}{YYYY}!'",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.ParseFlags(args)
			if err != nil {
				pkg.PrintFatal(err.Error())
			}

			if strings.ToLower(outputFormat) != "kerbrute" && strings.ToLower(outputFormat) != "netexec" {
				pkg.PrintFatal("Unknown outputFormat!")
			}

			if ldapPort == -1 {
				if ldapS {
					ldapPort = 636
				} else {
					ldapPort = 389
				}
			}

			pkg.RunLDAPQuery(ldapServer, ldapPort, ldapS, ntlm, username, password, hash, domain, ou, filter, outputFile, outputFormat, mask, pageSize, silent, cacheFile, noCache, forceRefresh)
		},
	}

	ldapServer               string
	ldapPort, pageSize       int
	ldapS, ntlm              bool
	username, password, hash string
	domain, ou, filter       string
	mask                     string
	outputFile               string
	outputFormat             string
	silent                   bool
	cacheFile                string
	noCache                  bool
	forceRefresh             bool
)

func init() {
	rootCmd.Flags().StringVarP(&ldapServer, "server", "s", "", "LDAP server address")
	rootCmd.Flags().IntVarP(&ldapPort, "port", "P", -1, "LDAP server port. Default: 389 for LDAP and 636 for LDAPS")
	rootCmd.Flags().IntVar(&pageSize, "pageSize", 500, "Page size")
	rootCmd.Flags().BoolVar(&ldapS, "ldaps", false, "LDAP over SSL/TLS")
	rootCmd.Flags().BoolVar(&ntlm, "ntlm", false, "Use NTLM authentication instead of basic LDAP authentication")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "Username. If no username is specified, an anonymous bind is attempted")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "Password. If no password is specified, an unauthenticated bind is attempted")
	rootCmd.Flags().StringVar(&hash, "hash", "", "NTLM Hash (Pass-the-Hash)")
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "FQDN")
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "(&(objectClass=User)(objectCategory=Person))", "LDAP Query Filter")
	rootCmd.Flags().StringVar(&ou, "ou", "", "Organizational Unit. E.g.: OU=Users,OU=GDATA")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file. Appends an incremental number if the file already exists")
	rootCmd.Flags().StringVar(&outputFormat, "outputformat", "kerbrute", "Output format. kerbrute creates a single file with user:pass, netexec creates two files, one with user and one with pass")
	rootCmd.Flags().StringVarP(&mask, "mask", "m", "", "Password mask. E.g.: Foobar{givenName#Reverse}{MonthGerman}{YYYY}!")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "Do not print the user attributes and the user:pass combos")
	rootCmd.PersistentFlags().StringVar(&cacheFile, "cache-file", "ldap_cache.json", "File to store cached LDAP data")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Disable caching of LDAP data")
	rootCmd.PersistentFlags().BoolVar(&forceRefresh, "force-refresh", false, "Force a new LDAP query and update the cache")
	rootCmd.MarkFlagRequired("server")
	rootCmd.MarkFlagRequired("domain")
	rootCmd.MarkFlagRequired("mask")
}

func getShortDescription() string {
	return "Query LDAP server to retrieve user information and generate a pw spraying list"
}

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

Mask Placeholder Modifiers
- #Reverse : Reverse the string
- #LeetBasic : Subsitute e:3, o:0, i:1, a:4
- #LeetBasicPlus : Subsitute e:3, o:0, i:1, a:@, t:7`
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
