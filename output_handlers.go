package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
)

// HTMLSeparator is a separator line in HTML

const (
	successString = "CODE:%d RespTime:%d %-3s - %3s \n" // Code, response time, URL, Title
	HTMLSeparator = "\n" + `<div class="separator"> separator </div>` + "\n"
	tableLine = "\n" + `<tr><td>{{.HTTPStatusCode}}</td><td>{{.URL}}<a href="{{.URL}}">Link</a></td><td>{{.Title}}</td></tr>` + "\n"
	htmlTemplate = `<html> <head> <style> body { background-color: #222222; color: #b7aa6b; font-family: 'Courier New', monospace; margin: 0; padding: 0; } .table-header, table { width: 100%; border-collapse: collapse; table-layout: fixed; } .table-header div, th, td { padding: 5px; text-align: center; /* Left alignment */ font-size: 14px; } tr { border-bottom: none; } tr:first-child { border-top: 1px solid #b7aa6b; } tr:last-child { border-bottom: 1px solid #b7aa6b; } tr:hover { background-color: #333333; } .separator { text-align: center; margin-top: 20px; margin-bottom: 20px; font-size: 16px; font-weight: bold; } </style> <meta charset="UTF-8"> </head> <body> <table>` + "\n" + `</table>` + "\n" + `</body>` + "\n" + `</html>` + "\n"
	openTableTag = `<table>` + "\n"
	closeTableTag = `</table>` + "\n"
	closeHtmlTag = `</html>` + "\n"
	closeBodyTag = `</body>` + "\n"
)
func WriteToStdOut(dmn Domain) {
	log.Printf(
		successString,
		dmn.HTTPStatusCode,
		dmn.ResponseTime.Milliseconds(),
		dmn.URL,
		dmn.Title)
	}

// Create report file if it doesn't exist
// Add separator if it exists
// FIXME: handle errors
func CheckReportFile() *os.File {
	// Check if report.html exists
	_, err := os.Stat("report.html")
	fileExists := !os.IsNotExist(err)

	// Open report.html in append mode
	f, err := os.OpenFile("report.html", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	// If file exists, remove the last line and append a separator
	if fileExists {
		err := truncateLastLine(f, 3)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		f.WriteString(HTMLSeparator)
		f.WriteString(openTableTag)
		f.WriteString(closeTableTag)
		f.WriteString(closeBodyTag)
		f.WriteString(closeHtmlTag)
	} else {
		// Write HTML header and table header if file doesn't exist
		f.WriteString(htmlTemplate)
	}

	return f
}


// truncateLastLine removes the last n lines from the file represented by the *os.File parameter.
func truncateLastLine(f *os.File, n int) error {
	// Move the file pointer to the beginning
	_, err := f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	// Read the file into a byte slice
	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Initialize a counter for newline characters
	newlineCount := 0

	// Find the appropriate number of newline characters and truncate the file
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '\n' {
			newlineCount++
			if newlineCount == n {
				err = f.Truncate(int64(i))
				if err != nil {
					return fmt.Errorf("error truncating file: %w", err)
				}
				break
			}
		}
	}

	// Move the file pointer back to the end for appending
	_, err = f.Seek(0, 2)
	if err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	return nil
}

func WriteToReport(f *os.File,dmn Domain) {
		// Prepare the HTML template
		tmpl := template.Must(template.New("status").Parse(tableLine))
	
		// Remove the last line from the file (the closing </table> tag and </html> tag)
		truncateLastLine(f, 4)
	
		// Check URLs and append to report.html
		err := tmpl.Execute(f, dmn)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return
		}

		// Add closing tags
		f.WriteString(closeTableTag + closeBodyTag + closeHtmlTag)
	
	}
