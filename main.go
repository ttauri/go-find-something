package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Config struct {
	charset       string
	domaindLength int
	tries         int
	zones         []string
}

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
func RunScan(ctx context.Context, logChan chan Domain, sch chan int, conf Config) {
	log.Print("Starting scan")
	sem := CreateSemaphore(2)

	domains := DomainList{}
	domains.generateDomainNamesList(conf)

	for i, domain := range domains.Domains {
		sch <- i
		sem.Acquire()
		go func(domain Domain) {
			domain.checkAvailability()
			logChan <- domain
			sem.Release()
		}(domain)

	}
}
func main() {
	conf := Config{
		charset:       "abcdefghijklmnopqrstuvwxyz1234567890",
		domaindLength: 7,
		tries:         100000,
		zones:         []string{".com", ".org", ".net"},
	}
	logChan := make(chan Domain)
	sitesChecked := make(chan int)
	// Do something with context
	ctx, _ := context.WithCancel(context.Background())
	go RunScan(ctx, logChan, sitesChecked, conf)
	go EventHandler(ctx, logChan)

	for {
		iteration := <-sitesChecked
		fmt.Printf("\rDomains checked: %d", iteration) // \r moves the cursor to the beginning of the line
		time.Sleep(100 * time.Millisecond)             // Sleep for a little bit to simulate work
	}

}
