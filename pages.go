package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Set up types
type Fn func()

func createHomePage(pages *tview.Pages, quitFunc Fn, appPrimitives AppPrimitives) *tview.List {
	homePage := tview.NewList()
	homePage.SetBorder(true)
	homePage.SetTitle("Main Menu")

	homePage.AddItem("Portolio", "Some explanatory text", 'a', func() {
		db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Get all accounts and add them up
		var account1 Account
		var account2 Account
		var account3 Account
		db.First(&account1, 1)
		db.First(&account2, 2)
		db.First(&account3, 3)

		updateAppCashCadBalances(appPrimitives, account1.CadBalance+account2.CadBalance+account3.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account1.UsdBalance+account2.UsdBalance+account3.UsdBalance)
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

func createPortfolioPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives) *tview.Grid {
	page := tview.NewGrid()
	page.SetBorder(true)
	page.SetTitle("Portfolio")
	page.SetRows(3, 0, 3)
	page.SetColumns(0, 0, 0)

	db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Set these up before the dropdown so you can use them in the selected func
	accountType := tview.NewTextView()
	accountType.SetText("Account: All")

	marketValue := tview.NewTextView()
	marketValue.SetText("Market Value: ")

	unrealizedPLValue := tview.NewTextView()
	unrealizedPLValue.SetDynamicColors(true)
	unrealizedPLValue.SetText("Unrealized P/L: xxxx")

	realizedPLValue := tview.NewTextView()
	realizedPLValue.SetDynamicColors(true)
	realizedPLValue.SetText("Realized P/L: xxxx")

	table := tview.NewTable()
	table.SetBorders(true)

	// Set up header columns
	table.SetCell(0, 0, tview.NewTableCell("STOCK"))
	table.SetCell(0, 1, tview.NewTableCell("PRICE"))
	table.SetCell(0, 2, tview.NewTableCell("POS"))
	table.SetCell(0, 3, tview.NewTableCell("P/L"))

	// Create rows with data
	for r := 1; r < 100; r++ {
		for c := 0; c < 4; c++ {
			table.SetCell(r, c, tview.NewTableCell("test"+" "+strconv.Itoa(r)))
		}
	}

	//===================================

	topFlex := tview.NewFlex()
	topFlex.SetDirection(tview.FlexColumn)

	dropDown := tview.NewDropDown()
	dropDown.SetBorder(true)

	dropDown.AddOption("All Accounts", func() {
		accountType.SetText("Account: All")
		marketValue.SetText("Market Value: 200000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]1000[white]")
		realizedPLValue.SetText("Realized P/L: [green]5000[white]")

		// Get all accounts and add them up
		var account1 Account
		var account2 Account
		var account3 Account
		db.First(&account1, 1)
		db.First(&account2, 2)
		db.First(&account3, 3)

		updateAppCashCadBalances(appPrimitives, account1.CadBalance+account2.CadBalance+account3.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account1.UsdBalance+account2.UsdBalance+account3.UsdBalance)

	})
	dropDown.AddOption("Margin", func() {
		accountType.SetText("Account: Margin")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]500[white]")
		realizedPLValue.SetText("Realized P/L: [green]3000[white]")

		// Get margin account
		var account Account
		db.First(&account, 1)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
	})
	dropDown.AddOption("TFSA", func() {
		accountType.SetText("Account: TFSA")
		marketValue.SetText("Market Value: 100000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]200[white]")
		realizedPLValue.SetText("Realized P/L: [green]1000[white]")

		// Get tfsa account
		var account Account
		db.First(&account, 2)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
	})
	dropDown.AddOption("RRSP", func() {
		accountType.SetText("Account: RRSP")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]300[white]")
		realizedPLValue.SetText("Realized P/L: [green]1000[white]")

		// Get tfsa account
		var account Account
		db.First(&account, 3)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
	})
	dropDown.SetCurrentOption(0)

	// Buttons
	buyButton := tview.NewButton("Buy").SetSelectedFunc(func() {
		pages.SwitchToPage("portfolioBuy")
	})

	sellButton := tview.NewButton("Sell").SetSelectedFunc(func() {
		pages.SwitchToPage("portfolioSell")
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
	leftMiddleFlex.SetTitle("Account Information")
	leftMiddleFlex.SetDirection(tview.FlexRow)

	cashBalancesFlex := tview.NewFlex()
	cashBalancesFlex.SetBorder(true)
	cashBalancesFlex.SetTitle("Cash Balances")
	cashBalancesFlex.SetDirection(tview.FlexRow)

	cashBalancesFlex.AddItem(appPrimitives.PortfolioCadTextView, 0, 1, false)
	cashBalancesFlex.AddItem(appPrimitives.PortfolioUsdTextView, 0, 1, false)

	leftMiddleFlex.AddItem(accountType, 0, 1, false)
	leftMiddleFlex.AddItem(marketValue, 0, 1, false)
	leftMiddleFlex.AddItem(realizedPLValue, 0, 1, false)
	leftMiddleFlex.AddItem(unrealizedPLValue, 0, 1, false)
	leftMiddleFlex.AddItem(cashBalancesFlex, 0, 2, false)

	rightMiddleFlex := tview.NewFlex()
	rightMiddleFlex.SetBorder(true)
	rightMiddleFlex.SetDirection(tview.FlexRow)

	rightMiddleFlex.AddItem(table, 0, 1, false)

	bottomFlex := tview.NewFlex()
	bottomFlex.SetBorder(true)
	bottomFlex.SetTitle("Navigation")
	bottomFlex.SetDirection(tview.FlexColumn)

	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](tab)[white] Switch between items"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](enter)[white] Select item"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](a)[white] Select account"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](b)[white] Buy stock"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](s)[white] Sell stock"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](t)[white] Select table"), 0, 1, false)
	bottomFlex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow](q)[white] Back"), 0, 1, false)

	// page.AddItem()
	page.AddItem(topFlex, 0, 0, 1, 3, 0, 0, true)
	page.AddItem(leftMiddleFlex, 1, 0, 1, 1, 0, 0, false)
	page.AddItem(rightMiddleFlex, 1, 1, 1, 2, 0, 0, false)
	page.AddItem(bottomFlex, 2, 0, 1, 3, 0, 0, false)

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
			} else if selected == 2 {
				selected = 3
				app.SetFocus(table)
			} else {
				selected = 0
				app.SetFocus(dropDown)
			}
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			dropDown.SetCurrentOption(0)
			pages.SwitchToPage("home")
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'a' {
			app.SetFocus(dropDown)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			app.SetFocus(buyButton)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 's' {
			app.SetFocus(sellButton)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			app.SetFocus(table)
		}

		return event

	})

	return page
}

func createPortfolioBuyPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle("Buy")

	page.AddButton("Buy", func() {

	})
	page.AddButton("Cancel", func() {
		pages.SwitchToPage("portfolio")
	})

	return page
}

func createPortfolioSellPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle("Sell")

	page.AddButton("Sell", func() {

	})
	page.AddButton("Cancel", func() {
		pages.SwitchToPage("portfolio")
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

func createFundsTransactionTypePage(pages *tview.Pages, appPrimitives AppPrimitives, title string) *tview.List {
	page := tview.NewList()
	page.SetBorder(true)
	page.SetTitle(fmt.Sprintf("Funds %s", title))

	db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if title == "Margin" {
		page.AddItem("Deposit", "Add funds to your margin account", 'a', func() { pages.SwitchToPage("fundsMarginDeposit") })
		page.AddItem("Withdraw", "Remove funds to your margin account", 'b', func() {

			var account Account
			db.First(&account, 1)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage("fundsMarginWithdraw")
		})

	} else if title == "TFSA" {
		page.AddItem("Deposit", "Add funds to your TFSA account", 'a', func() { pages.SwitchToPage("fundsTfsaDeposit") })
		page.AddItem("Withdraw", "Remove funds to your TFSA account", 'b', func() {

			var account Account
			db.First(&account, 2)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage("fundsTfsaWithdraw")
		})

	} else if title == "RRSP" {
		page.AddItem("Deposit", "Add funds to your RRSP account", 'a', func() { pages.SwitchToPage("fundsRrspDeposit") })
		page.AddItem("Withdraw", "Remove funds to your RRSP account", 'b', func() {

			var account Account
			db.First(&account, 3)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage("fundsRrspWithdraw")
		})
	}

	page.AddItem("Go back", "Press to go back", 'q', func() { pages.SwitchToPage("funds") })
	return page
}

func createFundsDepositPage(pages *tview.Pages, appPrimitives AppPrimitives, title string) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle(fmt.Sprintf("Deposit Funds %s", title))

	var selectedCurrency string = ""
	var selectedCurrencyId int

	page.AddInputField("Amount", "", 20, nil, nil)
	page.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })

	page.AddButton("Save", func() {
		db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		amountStr := page.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amount, err := strconv.ParseFloat(amountStr, 64)

		if selectedCurrency == "cad" {
			selectedCurrencyId = 1
		} else {
			selectedCurrencyId = 2
		}

		var accountId int

		// Set the proper account ID based on the title
		if title == "Margin" {
			accountId = 1
		} else if title == "TFSA" {
			accountId = 2
		} else if title == "RRSP" {
			accountId = 3
		}

		// Create transaction
		if err == nil {
			db.Create(&CashTransaction{Amount: amount, Type: "deposit", CurrencyID: selectedCurrencyId, AccountID: accountId})

			var account Account
			db.First(&account, accountId)

			// Update account balance for specific currency
			if selectedCurrency == "cad" {
				var newBalance float64 = account.CadBalance + amount
				account.CadBalance = newBalance
				updateAppCashCadBalances(appPrimitives, newBalance)
			} else {
				var newBalance float64 = account.UsdBalance + amount
				account.UsdBalance = newBalance
				updateAppCashUsdBalances(appPrimitives, newBalance)
			}

			// Save the updated record
			db.Save(&account)

		}

		page.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")

		if title == "Margin" {
			pages.SwitchToPage("fundsMargin")
		} else if title == "TFSA" {
			pages.SwitchToPage("fundsTfsa")
		} else if title == "RRSP" {
			pages.SwitchToPage("fundsRrsp")
		}

	})

	page.AddButton("Cancel", func() {
		page.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")

		if title == "Margin" {
			pages.SwitchToPage("fundsMargin")
		} else if title == "TFSA" {
			pages.SwitchToPage("fundsTfsa")
		} else if title == "RRSP" {
			pages.SwitchToPage("fundsRrsp")
		}
	})

	return page

}

func createFundsWithdrawPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives, title string) *tview.Flex {

	page := tview.NewFlex()
	page.SetBorder(true)
	page.SetTitle(fmt.Sprintf("Withdraw Funds %s", title))
	page.SetDirection(tview.FlexRow)

	balanceFlex := tview.NewFlex()
	balanceFlex.SetBorder(false)
	balanceFlex.SetDirection(tview.FlexRow)

	db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	balanceFlex.AddItem(appPrimitives.FundsCadTextView, 0, 1, false)
	balanceFlex.AddItem(appPrimitives.FundsUsdTextView, 0, 1, false)

	var selectedCurrency string = ""
	var selectedCurrencyId int

	form := tview.NewForm()
	form.SetBorder(true)

	form.AddInputField("Amount", "", 20, nil, nil)
	form.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })

	form.AddButton("Save", func() {

		amountStr := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amount, err := strconv.ParseFloat(amountStr, 64)

		if selectedCurrency == "cad" {
			selectedCurrencyId = 1
		} else {
			selectedCurrencyId = 2
		}

		var accountId int
		// Set the proper account ID based on the title
		if title == "Margin" {
			accountId = 1
		} else if title == "TFSA" {
			accountId = 2
		} else if title == "RRSP" {
			accountId = 3
		}

		// Create transaction
		if err == nil {
			db.Create(&CashTransaction{Amount: amount, Type: "withdraw", CurrencyID: selectedCurrencyId, AccountID: accountId})

			var account Account
			db.First(&account, accountId)

			var newBalance float64

			// Update account balance for specific currency
			if selectedCurrency == "cad" {
				newBalance = account.CadBalance - amount
				account.CadBalance = newBalance

				if newBalance >= 0 {
					updateAppCashCadBalances(appPrimitives, newBalance)
				}

			} else {
				newBalance = account.UsdBalance - amount
				account.UsdBalance = newBalance

				if newBalance >= 0 {
					updateAppCashUsdBalances(appPrimitives, newBalance)
				}
			}

			// Save the updated record
			if newBalance >= 0 {
				db.Save(&account)
			}

		}

		form.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")

		if title == "Margin" {
			pages.SwitchToPage("fundsMargin")
		} else if title == "TFSA" {
			pages.SwitchToPage("fundsTfsa")
		} else if title == "RRSP" {
			pages.SwitchToPage("fundsRrsp")
		}
	})

	form.AddButton("Cancel", func() {
		form.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")

		if title == "Margin" {
			pages.SwitchToPage("fundsMargin")
		} else if title == "TFSA" {
			pages.SwitchToPage("fundsTfsa")
		} else if title == "RRSP" {
			pages.SwitchToPage("fundsRrsp")
		}
	})

	page.AddItem(balanceFlex, 0, 1, false)
	page.AddItem(form, 0, 4, false)

	// Handle key presses -= Navigation
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(form)

		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {

			if title == "Margin" {
				pages.SwitchToPage("fundsMargin")
			} else if title == "TFSA" {
				pages.SwitchToPage("fundsTfsa")
			} else if title == "RRSP" {
				pages.SwitchToPage("fundsRrsp")
			}
		}

		return event

	})

	return page

}

func createStatisticsPage(pages *tview.Pages) *tview.TextView {
	statisticsPage := tview.NewTextView()
	statisticsPage.SetDynamicColors(true)
	statisticsPage.SetBorder(true)
	statisticsPage.SetTitle("Statistics")

	return statisticsPage
}
