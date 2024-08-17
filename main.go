package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Home page
	homePage := createHomePage(pages, func() { app.Stop() })
	// Portfolio Page
	portfolioPage := createPortfolioPage(pages, app)
	// Add Funds Page
	fundsPage := createFundsPage(pages)
	fundsMarginPage := createFundsMarginPage(pages)
	fundsTfsaPage := createFundsTfsaPage(pages)
	fundsRrspPage := createFundsRrspPage(pages)
	// Statisitcs Page
	statisticsPage := createStatisticsPage(pages)

	// Add pages to Pages
	pages.AddPage("home", homePage, true, true)
	pages.AddPage("portfolio", portfolioPage, true, false)
	pages.AddPage("funds", fundsPage, true, false)
	pages.AddPage("fundsMargin", fundsMarginPage, true, false)
	pages.AddPage("fundsTfsa", fundsTfsaPage, true, false)
	pages.AddPage("fundsRrsp", fundsRrspPage, true, false)
	pages.AddPage("statistics", statisticsPage, true, false)

	// Set the initial page to the home page
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}
