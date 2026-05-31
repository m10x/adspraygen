[![Release](https://img.shields.io/github/release/m10x/adspraygen.svg?color=brightgreen)](https://github.com/m10x/adspraygen/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/m10x/adspraygen)](https://goreportcard.com/report/github.com/m10x/adspraygen)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/m10x/adspraygen)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
# ADSprayGen

ADSprayGen a command-line utility written in Go that leverages LDAP (Lightweight Directory Access Protocol) to retrieve user attributes. These attributes can then be used to generate possible passwords for the users. A mask is required to generate the passwords, which can contain user attribute placeholders and modifiers for them.
- [Installation](#installation)
    - [Option 1: Prebuilt Binary](#option-1-prebuilt-binary)
    - [Option 2: Install Using Go](#option-2-install-using-go)
- [Usage](#usage)
    - [Mask Placeholders](#mask-placeholders)
    - [Modifiers](#modifiers)
- [TODOs](#todos)
- [Common LDAP Errors](#common-ldap-errors)


## Installation
### Option 1: Prebuilt Binary
Prebuilt binaries of ADSprayGen are provided on the [releases page](https://github.com/m10x/adspraygen/releases).
### Option 2: Install Using Go
Requirements: go1.23 or higher
```bash
go install -v github.com/m10x/adspraygen@latest
```

## Usage
ADSprayGen now provides 3 subcommands:

- `adspraygen pattern` - generic pattern generation from `[WORD]`, `[NUMBER]`, `[SPECIAL]`
  - use this one first to easily generate a *LOT* of password masks
- `adspraygen gen` - LDAP query and combo generation (previous default behavior).
  - use this one second to use masks to generate user:password combos.
  - queries LDAP and caches LDAP attributes to use them for password generation.
- `adspraygen spray` - kerbrute spray wrapper with lockout-safe waiting
  - use this one last to use kerbrute to spray the user:password combos.
  - uses the cached LDAP password policy information or via parameters specified password policy in order not to lock accounts

**Examples:**

- `adspraygen gen -d domain.local -u m10x -p m10x -s 10.10.10.10 -m 'Foobar{givenName#Reverse}{MonthGerman}{YYYY}!' -f '(&(objectClass=User)(objectCategory=Person))(!(userAccountControl:1.2.840.113556.1.4.803:=2))' -o spray.txt`
  - `-s`: DC IP
  - `-f`: LDAP Query, Here For Only Enabled Accounts
  - `-m`: The password mask to be used: alternatively use adspraygen pattern and adspraygen gen --mask-file in order to easily generate a lot of password masks
- `adspraygen pattern --patterns-file patterns.txt --nouns nouns.txt --out masks.txt --limit 50000`
- `adspraygen gen -d domain.local -u m10x -p m10x -s 10.10.10.10 --mask-file masks.txt -o spray.txt`
- `adspraygen spray --file spray.txt --domain domain.local --dc 10.10.10.10 --extra-flags "--safe"`
  - sprays the user:password combinations
  - uses the cached lockout policy information. Otherwise use --lockout-threshold and --reset-lockout-counter

### Mask Placeholders
- **{cn}** : Full Name
- **{givenName}** : First Name
- **{sn}** : Last Name
- **{sAMAccountName}** : Logon Name (Pre Windows 2000)
- **{rPrincipalName}** : Logon Name
- **{description}** : Description
- **{info}** : Notes
- **{department}** : Department
- **{I}** : City
- **{postcalCode}** : Postal Code
- Last password change
    - **{YYYY}** : e.g. 2024
    - **{YY}** : e.g. 24
    - **{MM}** : e.g. 01
    - **{M}** : e.g. 1
    - **{SeasonGerman}** : e.g. Herbst
    - **{SeasonAmerican}** : e.g. Fall
    - **{SeasonBritish}** : e.g. Autumn
    - **{MonthGerman}** : e.g. Januar
    - **{MonthEnglish}** : e.g. January

### Modifiers

Modifiers transform attribute values. Append with `#`, chain multiple modifiers with additional `#`.

| Modifier | Description | Example | Result |
|---|---|---|---|
| `#Upper` | Uppercase | `{givenName#Upper}` | `JOHN` |
| `#Lower` | Lowercase | `{givenName#Lower}` | `john` |
| `#Title` | Capitalize each word | `{givenName#Title}` | `John Smith` |
| `#Capitalize` | Capitalize first letter only | `{givenName#Capitalize}` | `John` |
| `#AlternateLower` | Alternating case, start lower | `{givenName#AlternateLower}` | `jOhN` |
| `#AlternateUpper` | Alternating case, start upper | `{givenName#AlternateUpper}` | `JoHn` |
| `#Reverse` | Reverse the string | `{givenName#Reverse}` | `nhoJ` |
| `#LeetBasic` | Substitute e→3, o→0, i→1, a→4 | `{givenName#LeetBasic}` | `J0hn` |
| `#LeetBasicPlus` | Like LeetBasic + a→@, t→7 | `{givenName#LeetBasicPlus}` | `J0hn` |
| `#Pattern(x>y)` | Replace x with y; chain rules with `;` | `{sn#Pattern(o>oO;a>4)}` | `JoOhn` |

#### Pattern Modifier Examples
```
{firstName#Pattern(a>4)}           // Replace all 'a' with '4'
{lastName#Pattern(e>3;a>4;i>1)}    // Leetspeak conversion
{department#Pattern(IT>tech)}      // Replace whole strings
{givenName#Pattern(ll>1)}          // Replace multiple characters
{userName#Pattern(a>)}             // Remove all 'a' characters
```

Escape special characters in patterns with a backslash:
```
{text#Pattern(a\;b>c)}     // Replaces "a;b" with "c"
{text#Pattern(a\>b>c)}     // Replaces "a>b" with "c"
{text#Pattern(\(x\)>y)}    // Replaces "(x)" with "y"
{text#Pattern(a\#b>c)}     // Replaces "a#b" with "c"
```

#### Chaining Modifiers
```
{firstName#Upper#Reverse}              // "john" → "NHOJ"
{department#Lower#Pattern(it>tech)}    // "IT Support" → "tech support"
{givenName#AlternateUpper#LeetBasic}   // "john" → "J0Hn"
```

## Common LDAP Errors
- `LDAP Result Code 1 "Operations Error": 000004DC: LdapErr: DSID-0C090A5C, comment: In order to perform this operation a successful bind must be completed on the connection.` - Anonymous/Unauthenticated bind is not possible. Specify a password or NTLM hash.
- `LDAP Result Code 49 "Invalid Credentials": 80090308: LdapErr: DSID-0C090439, comment: AcceptSecurityContext error` - The specified credentials are invalid
- `LDAP Result Code 49 "Invalid Credentials": 8009030C: LdapErr: DSID-0C0906B5, comment: AcceptSecurityContext error` - Unauthenticated NTLM bind is not possible or specified credentials are not valid.
