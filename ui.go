package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2/widget"
)

func CreateNewProgressBar(min float64, max float64) *widget.ProgressBar {
	progressBar := widget.NewProgressBar()
	progressBar.Min = min
	progressBar.Max = max
	return progressBar
}

func UpdateLogBox(ctx context.Context, logBox *widget.Entry, logChan <-chan Domain) {
	for {
		select {
		case domain, ok := <-logChan:
			if !ok {
				// logChan was closed
				return
			}
			if domain.Status != SiteDoesNotExist {
				fmt.Println("UpdateLogBox*****")
				logBox.CursorRow = len(logBox.Text) - 1
				// logBox.Append(fmt.Sprintf("%s %s %s\n", domain.Status, domain.URL, domain.Title))
				logBox.SetText(logBox.Text + fmt.Sprintf("%s %s %s\n", domain.Status, domain.URL, domain.Title))
			}
		case <-ctx.Done():
			logBox.CursorRow = len(logBox.Text) - 1
			logBox.SetText(logBox.Text + "\nScan was stopped\n")
			// Context was cancelled
			return
		}
	}
}

func UpdateProgressBar(ctx context.Context, progressBar *widget.ProgressBar, progressChan <-chan int) {
	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				// progressChan was closed
				return
			}
			progressBar.SetValue(float64(progress))
		case <-ctx.Done():
			// Context was cancelled
			return
		}
	}
}
