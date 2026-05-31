package cmd

import (
	"fmt"
	"strings"

	"github.com/m10x/adspraygen/pkg"

	"github.com/spf13/cobra"
)

var (
	ldapServer               string
	ldapPort, pageSize       int
	ldapS, ntlm              bool
	username, password, hash string
	domain, ou, filter       string
	mask                     string
	maskFile                 string
	outputFile               string
	outputFormat             string
	silent                   bool
	cacheFile                string
	noCache                  bool
	forceRefresh             bool
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Query LDAP and generate spray credentials",
	Long:    fmt.Sprintf("%s\n\n%s", getGenShortDescription(), getMaskOptions()),
	Example: "adspraygen gen -d domain.local -u m10x -p m10x -s 10.10.10.10 -m 'Foobar{givenName#Reverse}{MonthGerman}{YYYY}!'",
	Run: func(cmd *cobra.Command, args []string) {
		var masks []string
		if maskFile != "" {
			lines, err := pkg.ReadMaskFile(maskFile)
			if err != nil {
				pkg.PrintFatal(err.Error())
			}
			if len(lines) == 0 {
				pkg.PrintFatal("Mask file is empty")
			}
			masks = lines
		} else {
			masks = []string{mask}
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

		pkg.RunLDAPQuery(ldapServer, ldapPort, ldapS, ntlm, username, password, hash, domain, ou, filter, outputFile, outputFormat, masks, pageSize, silent, cacheFile, noCache, forceRefresh)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&ldapServer, "server", "s", "", "LDAP server address")
	genCmd.Flags().IntVarP(&ldapPort, "port", "P", -1, "LDAP server port. Default: 389 for LDAP and 636 for LDAPS")
	genCmd.Flags().IntVar(&pageSize, "pageSize", 500, "Page size")
	genCmd.Flags().BoolVar(&ldapS, "ldaps", false, "LDAP over SSL/TLS")
	genCmd.Flags().BoolVar(&ntlm, "ntlm", false, "Use NTLM authentication instead of basic LDAP authentication")
	genCmd.Flags().StringVarP(&username, "username", "u", "", "Username. If no username is specified, an anonymous bind is attempted")
	genCmd.Flags().StringVarP(&password, "password", "p", "", "Password. If no password is specified, an unauthenticated bind is attempted")
	genCmd.Flags().StringVar(&hash, "hash", "", "NTLM Hash (Pass-the-Hash)")
	genCmd.Flags().StringVarP(&domain, "domain", "d", "", "FQDN")
	genCmd.Flags().StringVarP(&filter, "filter", "f", "(&(objectClass=User)(objectCategory=Person))", "LDAP Query Filter")
	genCmd.Flags().StringVar(&ou, "ou", "", "Organizational Unit. E.g.: OU=Users,OU=GDATA")
	genCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file. Appends an incremental number if the file already exists")
	genCmd.Flags().StringVar(&outputFormat, "outputformat", "kerbrute", "Output format. kerbrute creates a single file with user:pass, netexec creates two files, one with user and one with pass")
	genCmd.Flags().StringVarP(&mask, "mask", "m", "", "Password mask. E.g.: Foobar{givenName#Reverse}{MonthGerman}{YYYY}!")
	genCmd.Flags().StringVar(&maskFile, "mask-file", "", "File with one mask per line (mutually exclusive with --mask)")
	genCmd.Flags().BoolVar(&silent, "silent", false, "Do not print the user attributes and the user:pass combos")
	genCmd.Flags().StringVar(&cacheFile, "cache-file", "ldap_cache.json", "File to store cached LDAP data")
	genCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching of LDAP data")
	genCmd.Flags().BoolVar(&forceRefresh, "force-refresh", false, "Force a new LDAP query and update the cache")

	genCmd.MarkFlagRequired("server")
	genCmd.MarkFlagRequired("domain")
	genCmd.MarkFlagsMutuallyExclusive("mask", "mask-file")
	genCmd.MarkFlagsOneRequired("mask", "mask-file")
}

func getGenShortDescription() string {
	return "Query LDAP server to retrieve user information and generate a pw spraying list"
}
