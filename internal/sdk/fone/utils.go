package fone

import (
	"log"
	"net"
	"regexp"
	"strings"
)

// IsOkDomain check domain
func IsOkDomain(domain string) bool {

	if domain == "" || len(strings.Replace(domain, ".", "", -1)) > 255 {
		return false
	}

	re := `^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*(\.[a-zA-Z0-9]{0,62}){1}\.?$`

	return match(re, domain)
}

// IsOkIP check ip
func IsOkIP(ip string) bool {

	return net.ParseIP(ip) != nil
}

func match(re, s string) bool {
	ok, err := regexp.MatchString(re, s)
	if err != nil {
		log.Println("re match error:", err)
		return false
	}
	if ok {
		return true
	}
	return false
}
