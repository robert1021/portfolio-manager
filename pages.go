package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Set up types
type Fn func()

func createHomePage(pages *tview.Pages, quitFunc Fn) *tview.List {
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

func createPortfolioPage(pages *tview.Pages, app *tview.Application) *tview.Grid {
	page := tview.NewGrid()
	page.SetBorder(true)
	page.SetTitle("Portfolio")
	page.SetRows(3, 0, 3)
	page.SetColumns(30, 0, 30)

	topFlex := tview.NewFlex()
	topFlex.SetDirection(tview.FlexColumn)

	dropDown := tview.NewDropDown()
	dropDown.SetBorder(true)
	dropDown.AddOption("All Accounts", func() { fmt.Println("all accounts") })
	dropDown.AddOption("Margin", func() { fmt.Println("margin") })
	dropDown.AddOption("TFSA", func() { fmt.Println("TFSA") })
	dropDown.AddOption("RRSP", func() { fmt.Println("RRSP") })

	buyButton := tview.NewButton("Buy").SetSelectedFunc(func() {
		fmt.Println("Buy")
	})

	topFlex.AddItem(dropDown, 0, 1, true)
	topFlex.AddItem(buyButton, 0, 1, false)

	// Create boxes for the other rows
	box1 := tview.NewBox().SetBorder(true).SetTitle("Row 2")
	box2 := tview.NewBox().SetBorder(true).SetTitle("Row 3")
	box3 := tview.NewBox().SetBorder(true).SetTitle("Row 3")
	box4 := tview.NewBox().SetBorder(true).SetTitle("Row 3")

	// page.AddItem()
	page.AddItem(topFlex, 0, 0, 1, 3, 0, 0, true)
	page.AddItem(box1, 1, 0, 1, 1, 0, 100, false)
	page.AddItem(box2, 1, 1, 1, 1, 0, 100, false)
	page.AddItem(box3, 1, 2, 1, 1, 0, 100, false)
	page.AddItem(box4, 2, 0, 1, 3, 0, 0, false)

	var selected int8

	// Handle key presses -= Navigation
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if selected == 0 {
				selected = 1
				app.SetFocus(buyButton)
			} else {
				selected = 0
				app.SetFocus(dropDown)
			}
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			pages.SwitchToPage("home")
		}

		return event

	})

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
