package main

import (
	"github.com/rivo/tview"
)

type fn func()

func createHomePage(pages *tview.Pages, quitFunc fn) *tview.List {
	homePage := tview.NewList()
	homePage.SetBorder(true)
	homePage.SetTitle("Main Menu")

	homePage.AddItem("Portolio", "Some explanatory text", 'a', func() {
		pages.SwitchToPage("portfolio")
	})
	homePage.AddItem("Funds", "Some explanatory text", 'b', func() {
		pages.SwitchToPage("funds")
	})
	homePage.AddItem("Statistics", "Some explanatory text", 'c', func() {
		pages.SwitchToPage("statistics")
	})
	homePage.AddItem("Quit", "Press to exit", 'q', func() {
		quitFunc()
	})

	return homePage
}

func createPortfolioPage(pages *tview.Pages) *tview.TextView {
	page := tview.NewTextView()
	page.SetBorder(true)
	page.SetTitle("Portfolio")

	return page
}

func createFundsPage(pages *tview.Pages) *tview.List {
	fundsPage := tview.NewList()
	fundsPage.SetBorder(true)
	fundsPage.SetTitle("Funds")

	fundsPage.AddItem("Margin", "Some explanatory text", 'a', func() {
		pages.SwitchToPage("fundsMargin")
	})
	fundsPage.AddItem("TFSA", "Some explanatory text", 'b', func() {
		pages.SwitchToPage("fundsTfsa")
	})
	fundsPage.AddItem("RRSP", "Some explanatory text", 'c', func() {
		pages.SwitchToPage("fundsRrsp")
	})
	fundsPage.AddItem("Go back", "Press to go back", 'q', func() {
		pages.SwitchToPage("home")
	})

	return fundsPage
}

func createFundsMarginPage(pages *tview.Pages) *tview.List {
	page := tview.NewList()
	page.SetBorder(true)
	page.SetTitle("Funds Margin")

	page.AddItem("Go back", "Press to go back", 'q', func() { pages.SwitchToPage("funds") })

	return page

}

func createFundsTfsaPage(pages *tview.Pages) *tview.List {
	page := tview.NewList()
	page.SetBorder(true)
	page.SetTitle("Funds TFSA")

	page.AddItem("Go back", "Press to go back", 'q', func() { pages.SwitchToPage("funds") })

	return page

}

func createFundsRrspPage(pages *tview.Pages) *tview.List {
	page := tview.NewList()
	page.SetBorder(true)
	page.SetTitle("Funds RRSP")

	page.AddItem("Go back", "Press to go back", 'q', func() { pages.SwitchToPage("funds") })

	return page

}

func createStatisticsPage(pages *tview.Pages) *tview.TextView {
	statisticsPage := tview.NewTextView()
	statisticsPage.SetDynamicColors(true)
	statisticsPage.SetBorder(true)
	statisticsPage.SetTitle("Statistics")

	return statisticsPage
}
