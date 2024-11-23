package common

import (
	"os"
	"strconv"
)

// ParseUint safely parses a string to uint64
func ParseUint(s string) uint64 {
	value, _ := strconv.ParseUint(s, 10, 64)
	return value

}

// hostname needed to append into metrics
func GetHostname() string {

	// HOSTNAME env variable takes precedence
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	return hostname
}
