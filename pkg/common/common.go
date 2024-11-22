package common

import (
	"log"
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
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v\n", err)
		return ""
	}
	return hostname
}
