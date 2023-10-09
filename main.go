package main

import (
	"context"
	"fmt"
	"time"
)

// "fyne.io/fyne/v2"
// "fyne.io/fyne/v2/app"
// "fyne.io/fyne/v2/container"
// "fyne.io/fyne/v2/theme"
// "fyne.io/fyne/v2/widget"

type Config struct {
	charset       string
	domaindLength int
	tries         int
	zones         []string
}

// type CustomTheme struct {
// 	fyne.Theme
// }

// func (m *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
// 	if name == theme.ColorNameDisabled {
// 		return color.White // Change this to any color you prefer
// 	}
// 	return m.Theme.Color(name, variant)
// }

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
		time.Sleep(100 * time.Millisecond)        // Sleep for a little bit to simulate work
	}

	// myApp := app.NewWithID("Go find something")
	// // TODO: Try to change theme
	// myApp.Settings().SetTheme(&CustomTheme{Theme: theme.DefaultTheme()})
	// myWindow := myApp.NewWindow("App with Log Box")

	// progressBar := CreateNewProgressBar(0, 10000)

	// // Log box
	// sitesLogBox := widget.NewMultiLineEntry()
	// sitesLogBox.Disable()
	// sitesLogBox.SetText("Sites log box\ndfjdhfj")

	// generalLogBox := widget.NewMultiLineEntry()
	// generalLogBox.Disable()

	// progressChan := make(chan int)
	// logChan := make(chan Domain)

	// ctx, cancel := context.WithCancel(context.Background())

	// startButton := widget.NewButton("Scan", func() {
	// 	go RunScan(ctx, progressChan, logChan)
	// 	go UpdateLogBox(ctx, generalLogBox, logChan)
	// 	go UpdateProgressBar(ctx, progressBar, progressChan)
	// })

	// stopButton := widget.NewButton("Stop", func() {
	// 	cancel()
	// })

	// // Row with start and stop buttons
	// firstRow := container.NewHBox(startButton, stopButton)
	// secondRow := container.NewVBox(progressBar)

	// // Tabs with log messages
	// tabs := container.NewAppTabs(
	// 	container.NewTabItem("Log", generalLogBox),
	// 	container.NewTabItem("Sites", sitesLogBox),
	// )

	// // Part of the UI that will displays all controls and progress bar
	// rowsWithControls := container.NewGridWithRows(2, firstRow, secondRow)
	// content := container.NewBorder(rowsWithControls, nil, nil, nil, tabs)

	// // Run app
	// myWindow.SetContent(content)
	// myWindow.Resize(fyne.NewSize(400, 300))
	// myWindow.CenterOnScreen()
	// myWindow.ShowAndRun()

}
