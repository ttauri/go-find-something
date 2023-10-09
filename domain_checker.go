package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type SiteStatus string

const (
	SiteExists       SiteStatus = "Exists"
	SiteDoesNotExist SiteStatus = "Does Not Exist"
	SiteUnavailable  SiteStatus = "Unavailable"
)

type Domain struct {
	URL            string
	Status         SiteStatus
	HTTPStatusCode int
	ResponseTime   time.Duration
	Error          error
	ResolveMethod  string
	Title          string
	AdditionalInfo map[string]interface{}
}

type DomainList struct {
	Domains []Domain
	Length  int
}
type Semaphore interface {
	Acquire()
	Release()
}

type semaphore struct {
	semChan chan struct{}
}

func (s *semaphore) Acquire() {
	s.semChan <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.semChan
}

func CreateSemaphore(maxRoutines int) Semaphore {
	return &semaphore{
		semChan: make(chan struct{}, maxRoutines),
	}
}

func (d *DomainList) deduplicateDomains() {
	seen := make(map[string]bool)
	var unique []Domain

	for _, domain := range d.Domains {
		if _, exists := seen[domain.URL]; !exists {
			seen[domain.URL] = true
			unique = append(unique, domain)
		}
	}
	d.Domains = unique
}

// Generate list of random domain names
func (d *DomainList) generateDomainNamesList(conf Config) {
	domains := []Domain{}
	for i := 0; i < conf.tries; i++ {
		domains = append(domains, generateRandomDomainName(conf))
	}
	d.Domains = domains
	d.deduplicateDomains()
}

// Get title from HTML page
func getTitle(resp *http.Response) string {
	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil { // Add this nil check
			return
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil { // Add this nil check
				title = n.FirstChild.Data
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return strings.TrimSpace(title)
}

// Check if site is available
func (domain *Domain) checkAvailability() {
	t1 := time.Now()
	client := GetHttpClient()
	resp, err := client.Get(domain.URL)
	respTime := time.Since(t1)

	if err != nil {
		if _, ok := err.(*url.Error); ok {
			domain.Status = SiteDoesNotExist
			return
		} else {
			domain.Status = SiteUnavailable
			return
		}
	}

	if resp.StatusCode == http.StatusOK {
		domain.Status = SiteExists
		domain.ResponseTime = respTime
		domain.Title = getTitle(resp)
	} else {
		domain.Status = SiteUnavailable
	}
	domain.HTTPStatusCode = resp.StatusCode
	defer resp.Body.Close()
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
