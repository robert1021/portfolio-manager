package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Home page
	homePage := tview.NewList()
	homePage.SetBorder(true)
	homePage.SetTitle("Main Menu")

	homePage.AddItem("Portolio", "Some explanatory text", 'a', func() {
		pages.SwitchToPage("portfolio")
	})
	homePage.AddItem("Add Funds", "Some explanatory text", 'b', func() {
		pages.SwitchToPage("addfunds")
	})
	homePage.AddItem("Statistics", "Some explanatory text", 'c', func() {
		pages.SwitchToPage("statistics")
	})
	homePage.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	// Portfolio Page
	portfolioPage := tview.NewTextView().
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Portfolio")

	// Add Funds Page
	addFundsPage := tview.NewTextView().
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Add Funds")

	// Statisitcs Page
	statisticsPage := tview.NewTextView().
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Statistics")

	// Add pages to Pages
	pages.AddPage("home", homePage, true, true)
	pages.AddPage("portfolio", portfolioPage, true, false)
	pages.AddPage("addfunds", addFundsPage, true, false)
	pages.AddPage("statistics", statisticsPage, true, false)

	// Set the initial page to the home page
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}
