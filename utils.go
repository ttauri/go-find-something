package main

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net/http"
	"time"
)



func EventHandler(ctx context.Context, logChan <-chan Domain) {
	reportFile := CheckReportFile()
	for {
		select {
		case site, ok := <-logChan:
			if !ok {
				// logChan was closed
				reportFile.Close()
				return
			}
			if site.Status != SiteDoesNotExist {
				WriteToReport(reportFile, site)
				WriteToStdOut(site)
			}
		case <-ctx.Done():
			// Context was cancelled
				reportFile.Close()
			return
		}
	}
}


// Generate random string with given length
func GenerateRandomString(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}



// Get HTTP client with disabled TLS verification
func GetHttpClient() *http.Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return httpClient
}

