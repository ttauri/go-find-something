package main

import (
	"math/rand"
	"time"
)

// Generate random string with given length
func GenerateRandomString(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// Generate random domain name with given length
func GenerateRandomDomainName(conf Config) Domain {
	url := "https://" + GenerateRandomString(conf.domaindLength, conf.charset) + conf.zones[rand.Intn(len(conf.zones))]
	domain := Domain{URL: url}
	return domain
}

// Generate random domain name with given length
func generateRandomDomainName(conf Config) Domain {
	url := "https://" + GenerateRandomString(conf.domaindLength, conf.charset) + conf.zones[rand.Intn(len(conf.zones))]
	domain := Domain{URL: url}
	return domain
}
