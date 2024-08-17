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
	page.SetColumns(0, 0, 0)

	// Set these up before the dropdown so you can use them in the selected func
	accountType := tview.NewTextView()
	accountType.SetText("Account: All")

	marketValue := tview.NewTextView()
	marketValue.SetText("Market Value: ")

	unrealizedPLValue := tview.NewTextView()
	unrealizedPLValue.SetText("Unrealized P/L: xxxx")

	realizedPLValue := tview.NewTextView()
	realizedPLValue.SetText("Realized P/L: xxxx")

	cadCashValue := tview.NewTextView()
	cadCashValue.SetText("CAD Cash: xxxx")

	usdCashValue := tview.NewTextView()
	usdCashValue.SetText("USD Cash: xxxx")

	//===================================

	topFlex := tview.NewFlex()
	topFlex.SetDirection(tview.FlexColumn)

	dropDown := tview.NewDropDown()
	dropDown.SetBorder(true)

	dropDown.AddOption("All Accounts", func() {
		accountType.SetText("Account: All")
		marketValue.SetText("Market Value: 200000")
		unrealizedPLValue.SetText("Unrealized P/L: 1000")
		realizedPLValue.SetText("Realized P/L: 5000")
		cadCashValue.SetText("CAD Cash: 800")
		usdCashValue.SetText("USD Cash: 800")
	})
	dropDown.AddOption("Margin", func() {
		accountType.SetText("Account: Margin")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: 500")
		realizedPLValue.SetText("Realized P/L: 3000")
		cadCashValue.SetText("CAD Cash: 500")
		usdCashValue.SetText("USD Cash: 500")
	})
	dropDown.AddOption("TFSA", func() {
		accountType.SetText("Account: TFSA")
		marketValue.SetText("Market Value: 100000")
		unrealizedPLValue.SetText("Unrealized P/L: 200")
		realizedPLValue.SetText("Realized P/L: 1000")
		cadCashValue.SetText("CAD Cash: 200")
		usdCashValue.SetText("USD Cash: 200")
	})
	dropDown.AddOption("RRSP", func() {
		accountType.SetText("Account: RRSP")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: 300")
		realizedPLValue.SetText("Realized P/L: 1000")
		cadCashValue.SetText("CAD Cash: 100")
		usdCashValue.SetText("USD Cash: 100")
	})
	dropDown.SetCurrentOption(0)

	// Buttons
	buyButton := tview.NewButton("Buy").SetSelectedFunc(func() {
		fmt.Println("Buy")
	})

	sellButton := tview.NewButton("Sell").SetSelectedFunc(func() {
		fmt.Println("Sell")
	})

	topFlex.AddItem(dropDown, 0, 2, true)
	topFlex.AddItem(tview.NewBox(), 20, 0, false)
	topFlex.AddItem(buyButton, 0, 1, false)
	topFlex.AddItem(tview.NewBox(), 5, 0, false)
	topFlex.AddItem(sellButton, 0, 1, false)
	topFlex.AddItem(tview.NewBox(), 1, 0, false)

	// Left middle flex layout
	leftMiddleFlex := tview.NewFlex()
	leftMiddleFlex.SetBorder(true)
	leftMiddleFlex.SetTitle("2")
	leftMiddleFlex.SetDirection(tview.FlexRow)

	cashBalancesFlex := tview.NewFlex()
	cashBalancesFlex.SetBorder(true)
	cashBalancesFlex.SetTitle("Cash Balances")
	cashBalancesFlex.SetDirection(tview.FlexRow)

	cashBalancesFlex.AddItem(cadCashValue, 0, 1, false)
	cashBalancesFlex.AddItem(usdCashValue, 0, 1, false)

	leftMiddleFlex.AddItem(accountType, 0, 1, false)
	leftMiddleFlex.AddItem(marketValue, 0, 1, false)
	leftMiddleFlex.AddItem(realizedPLValue, 0, 1, false)
	leftMiddleFlex.AddItem(unrealizedPLValue, 0, 1, false)
	leftMiddleFlex.AddItem(cashBalancesFlex, 0, 2, false)

	box2 := tview.NewBox().SetBorder(true).SetTitle("3")
	box4 := tview.NewBox().SetBorder(true).SetTitle("4")

	// page.AddItem()
	page.AddItem(topFlex, 0, 0, 1, 3, 0, 0, true)
	page.AddItem(leftMiddleFlex, 1, 0, 1, 1, 0, 0, false)
	page.AddItem(box2, 1, 1, 1, 2, 0, 0, false)
	page.AddItem(box4, 2, 0, 1, 3, 0, 0, false)

	var selected int8

	// Handle key presses -= Navigation
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if selected == 0 {
				selected = 1
				app.SetFocus(buyButton)
			} else if selected == 1 {
				selected = 2
				app.SetFocus(sellButton)
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
