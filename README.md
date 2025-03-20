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
Requirements: go1.21 or higher
```bash
go install -v github.com/m10x/adspraygen@latest
```

## Usage
Example: `adspraygen -d domain.local -u m10x -p m10x -s 10.10.10.10 -m 'Foobar{givenName#Reverse}{MonthGerman}{YYYY}!'`

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

Modifiers can be used to transform attribute values. Multiple modifiers can be chained using `#`.

#### Case Modifiers
- `#Upper` - Converts text to uppercase
  - Example: `{firstName#Upper}` → "JOHN"
- `#Lower` - Converts text to lowercase
  - Example: `{firstName#Lower}` → "john"
- `#Title` - Capitalizes first letter of each word
  - Example: `{firstName#Title}` → "John Smith"
- `#Capitalize` - Capitalizes only the first letter
  - Example: `{firstName#Capitalize}` → "John"


#### Alternating Case Modifiers
- `#AlternateLower` - Alternates case starting with lowercase
  - Example: `{firstName#Alternate}` → "jOhN"
- `#AlternateUpper` - Alternates case starting with uppercase
  - Example: `{firstName#AlternateUpper}` → "JoHn"

#### Text Transformation
- `#Reverse` - Reverses the text
  - Example: `{firstName#Reverse}` → "nhoJ"
- `#Pattern(from>to)` - Replaces characters according to pattern rules
  - Example: `{firstName#Pattern(a>4)}` → Replaces all 'a' with '4'
  - Example: `{lastName#Pattern(e>3;a>4;i>1)}` → Leetspeak conversion
  - Example: `{department#Pattern(IT>tech)}` → Replaces whole strings
  - Example: `{givenName#Pattern(ll>1)}` → Replaces multiple characters
  - Example: `{userName#Pattern(a>)}` → Removes all 'a' characters

#### Leetspeak Modifiers
- `#LeetBasic` - Basic leetspeak conversion (A->4, E->3, I->1, O->0)
  - Example: `{firstName#LeetBasic}` → "J0hn"
- `#LeetBasicPlus` - Extended leetspeak (Basic + A->@, T->7)
  - Example: `{firstName#LeetBasicPlus}` → "J0hn"

#### Chaining Modifiers
You can combine multiple modifiers to achieve complex transformations:
```
{firstName#Upper#Reverse}              // "JOHN" → "NHOJ"
{department#Lower#Pattern(it>tech)}    // "IT Support" → "tech support"
{userName#Camel#LeetBasic}            // "john doe" → "j0hnD03"
{text#AlternateUpper#Pattern(O>0)}    // "hello" → "H3Ll0"
```

Note: When using special characters in patterns (;, >, (, ), #), you need to escape them with a backslash:
```
{text#Pattern(a\;b>c)}     // Replaces "a;b" with "c"
{text#Pattern(a\>b>c)}     // Replaces "a>b" with "c"
{text#Pattern(\(x\)>y)}    // Replaces "(x)" with "y"
{text#Pattern(a\#b>c)}     // Replaces "a#b" with "c"
```

## TODOs
- dump LDAP user attributes
- handling of unknown mask attribute and unknown mask transformator
- dump Hostnames
- handling of givenName with multiple names

## Common LDAP Errors
- `LDAP Result Code 1 "Operations Error": 000004DC: LdapErr: DSID-0C090A5C, comment: In order to perform this operation a successful bind must be completed on the connection.` - Anonymous/Unauthenticated bind is not possible. Specify a password or NTLM hash.
- `LDAP Result Code 49 "Invalid Credentials": 80090308: LdapErr: DSID-0C090439, comment: AcceptSecurityContext error` - The specified credentials are invalid
- `LDAP Result Code 49 "Invalid Credentials": 8009030C: LdapErr: DSID-0C0906B5, comment: AcceptSecurityContext error` - Unauthenticated NTLM bind is not possible or specified credentials are not valid.