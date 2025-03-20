package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// CachedLDAPData represents the cached LDAP data
type CachedLDAPData struct {
	Entries    []LDAPEntry `json:"entries"`
	CachedAt   time.Time   `json:"cached_at"`
	SearchBase string      `json:"search_base"`
	LDAPFilter string      `json:"ldap_filter"`
	Attributes []string    `json:"attributes"`
	LDAPServer string      `json:"ldap_server"`
	LDAPPort   int         `json:"ldap_port"`
}

// LDAPEntry represents a single LDAP entry
type LDAPEntry struct {
	DN         string              `json:"dn"`
	Attributes map[string][]string `json:"attributes"`
}

// SaveLDAPDataToCache stores the LDAP data in a JSON file
func SaveLDAPDataToCache(entries []*ldap.Entry, searchBase, ldapFilter string, attributes []string, ldapServer string, ldapPort int, cacheFile string) error {
	cachedData := CachedLDAPData{
		CachedAt:   time.Now(),
		SearchBase: searchBase,
		LDAPFilter: ldapFilter,
		Attributes: attributes,
		LDAPServer: ldapServer,
		LDAPPort:   ldapPort,
	}

	// Convert LDAP entries to cache format
	for _, entry := range entries {
		cacheEntry := LDAPEntry{
			DN:         entry.DN,
			Attributes: make(map[string][]string),
		}

		for _, attr := range entry.Attributes {
			cacheEntry.Attributes[attr.Name] = attr.Values
		}

		cachedData.Entries = append(cachedData.Entries, cacheEntry)
	}

	// Save as JSON
	jsonData, err := json.MarshalIndent(cachedData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}

	err = os.WriteFile(cacheFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error saving cache file: %v", err)
	}

	return nil
}

// LoadLDAPDataFromCache loads the LDAP data from a JSON file
func LoadLDAPDataFromCache(cacheFile string) (*CachedLDAPData, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %v", err)
	}

	var cachedData CachedLDAPData
	err = json.Unmarshal(data, &cachedData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}

	return &cachedData, nil
}

// ConvertCacheToLDAPEntries converts cache entries back to LDAP entries
func ConvertCacheToLDAPEntries(cachedData *CachedLDAPData) []*ldap.Entry {
	var entries []*ldap.Entry

	for _, cacheEntry := range cachedData.Entries {
		entry := &ldap.Entry{
			DN: cacheEntry.DN,
		}

		for name, values := range cacheEntry.Attributes {
			entry.Attributes = append(entry.Attributes, &ldap.EntryAttribute{
				Name:   name,
				Values: values,
			})
		}

		entries = append(entries, entry)
	}

	return entries
}

// ShouldUpdateCache determines if the cache should be updated based on server information
func ShouldUpdateCache(cachedData *CachedLDAPData, ldapServer string, ldapPort int) bool {
	return cachedData.LDAPServer != ldapServer || cachedData.LDAPPort != ldapPort
}
