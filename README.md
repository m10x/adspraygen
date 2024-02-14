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
    - [Mask Placeholder Transformators](#mask-placeholder-transformators)
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

### Mask Placeholder Transformators
- **\#Reverse** : Reverse the string
- **\#LeetBasic** : Subsitute e:3, o:0, i:1, a:4
- **\#LeetBasicPlus** : Subsitute e:3, o:0, i:1, a:@, t:7

## TODOs
- dump LDAP user attributes
- import dumped LDAP user attributes instead of querying the LDAP server

## Common LDAP Errors
`LDAP Result Code 1 "Operations Error": 000004DC: LdapErr: DSID-0C090A5C, comment: In order to perform this operation a successful bind must be completed on the connection.`: Anonymous/Unauthenticated bind is not possible. Specify a password or NTLM hash.
`LDAP Result Code 49 "Invalid Credentials": 80090308: LdapErr: DSID-0C090439, comment: AcceptSecurityContext error` - The specified credentials are invalid
`LDAP Result Code 49 "Invalid Credentials": 8009030C: LdapErr: DSID-0C0906B5, comment: AcceptSecurityContext error` - Unauthenticated NTLM bind is not possible or specified credentials are not valid.
