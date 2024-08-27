package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Set up types
type Fn func()

func createHomePage(pages *tview.Pages, quitFunc Fn, appPrimitives AppPrimitives) *tview.List {
	homePage := tview.NewList()
	homePage.SetBorder(true)
	homePage.SetTitle("Main Menu")

	homePage.AddItem("Portolio", "Some explanatory text", 'a', func() {
		db := connectDb()

		// Get all accounts and add them up
		var account1 Account
		var account2 Account
		var account3 Account
		db.First(&account1, 1)
		db.First(&account2, 2)
		db.First(&account3, 3)

		updateAppCashCadBalances(appPrimitives, account1.CadBalance+account2.CadBalance+account3.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account1.UsdBalance+account2.UsdBalance+account3.UsdBalance)
		pages.SwitchToPage(portfolioPageName)
	})
	homePage.AddItem("Funds", "Some explanatory text", 'b', func() {
		pages.SwitchToPage(fundsPageName)
	})
	homePage.AddItem("Statistics", "Some explanatory text", 'c', func() {
		pages.SwitchToPage(statisticsPageName)
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

	db := connectDb()

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

	appPrimitives.PortfolioStockTable.SetBorders(true)

	// Create inital table
	updatePortfolioStockTable(queryAllStocks(db), appPrimitives)

	topFlex := tview.NewFlex()
	topFlex.SetDirection(tview.FlexColumn)

	// TODO:
	appPrimitives.PortfolioDropdown.SetBorder(true)

	appPrimitives.PortfolioDropdown.AddOption("All Accounts", func() {
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
		updatePortfolioStockTable(queryAllStocks(db), appPrimitives)

	})
	appPrimitives.PortfolioDropdown.AddOption("Margin", func() {
		accountType.SetText("Account: Margin")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]500[white]")
		realizedPLValue.SetText("Realized P/L: [green]3000[white]")

		// Get margin account
		var account Account
		db.First(&account, 1)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
		updatePortfolioStockTable(queryStocks(db, 1), appPrimitives)
	})
	appPrimitives.PortfolioDropdown.AddOption("TFSA", func() {
		accountType.SetText("Account: TFSA")
		marketValue.SetText("Market Value: 100000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]200[white]")
		realizedPLValue.SetText("Realized P/L: [green]1000[white]")

		// Get tfsa account
		var account Account
		db.First(&account, 2)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
		updatePortfolioStockTable(queryStocks(db, 2), appPrimitives)
	})
	appPrimitives.PortfolioDropdown.AddOption("RRSP", func() {
		accountType.SetText("Account: RRSP")
		marketValue.SetText("Market Value: 50000")
		unrealizedPLValue.SetText("Unrealized P/L: [red]300[white]")
		realizedPLValue.SetText("Realized P/L: [green]1000[white]")

		// Get tfsa account
		var account Account
		db.First(&account, 3)

		updateAppCashCadBalances(appPrimitives, account.CadBalance)
		updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
		updatePortfolioStockTable(queryStocks(db, 3), appPrimitives)
	})
	appPrimitives.PortfolioDropdown.SetCurrentOption(0)

	// Buttons
	buyButton := tview.NewButton("Buy").SetSelectedFunc(func() {
		pages.SwitchToPage(portfolioBuyPageName)
	})

	sellButton := tview.NewButton("Sell").SetSelectedFunc(func() {
		pages.SwitchToPage(portfolioSellPageName)
	})

	topFlex.AddItem(appPrimitives.PortfolioDropdown, 0, 2, true)
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

	rightMiddleFlex.AddItem(appPrimitives.PortfolioStockTable, 0, 1, false)

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
				app.SetFocus(appPrimitives.PortfolioStockTable)
			} else {
				selected = 0
				app.SetFocus(appPrimitives.PortfolioDropdown)
			}
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			appPrimitives.PortfolioDropdown.SetCurrentOption(0)
			pages.SwitchToPage(homePageName)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'a' {
			app.SetFocus(appPrimitives.PortfolioDropdown)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			app.SetFocus(buyButton)
			pages.SwitchToPage(portfolioBuyPageName)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 's' {
			app.SetFocus(sellButton)
			pages.SwitchToPage(portfolioSellPageName)
		} else if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			app.SetFocus(appPrimitives.PortfolioStockTable)
		}

		return event

	})

	return page
}

func createPortfolioBuyPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle("Buy")

	var selectedCurrency string = ""
	var selectedAccount string = ""
	var costStr string = ""

	page.AddInputField("Stock", "", 20, nil, nil)
	page.AddDropDown("Account", []string{"Margin", "TFSA", "RRSP"}, 0, func(option string, optionIndex int) { selectedAccount = option })
	page.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })
	page.AddInputField("Price", "", 20, nil, func(text string) {

		priceStr := page.GetFormItemByLabel("Price").(*tview.InputField).GetText()
		price, priceErr := strconv.ParseFloat(priceStr, 64)

		sharesStr := page.GetFormItemByLabel("Quantity").(*tview.InputField).GetText()
		shares, sharesErr := strconv.ParseFloat(sharesStr, 64)

		if priceErr == nil && sharesErr == nil {
			var amount float64 = shares * price
			costStr = fmt.Sprintf("%f", amount)
			page.GetFormItemByLabel("Cost").(*tview.TextView).SetText(costStr)

		} else {
			costStr = ""
			page.GetFormItemByLabel("Cost").(*tview.TextView).SetText(costStr)
		}

	})
	page.AddInputField("Quantity", "", 20, nil, func(text string) {

		priceStr := page.GetFormItemByLabel("Price").(*tview.InputField).GetText()
		price, priceErr := strconv.ParseFloat(priceStr, 64)

		sharesStr := page.GetFormItemByLabel("Quantity").(*tview.InputField).GetText()
		shares, sharesErr := strconv.ParseFloat(sharesStr, 64)

		if priceErr == nil && sharesErr == nil {
			var amount float64 = shares * price
			costStr = fmt.Sprintf("%f", amount)
			page.GetFormItemByLabel("Cost").(*tview.TextView).SetText(costStr)
		} else {
			costStr = ""
			page.GetFormItemByLabel("Cost").(*tview.TextView).SetText(costStr)
		}

	})
	page.AddTextView("Cost", "", 20, 2, true, false)

	page.AddButton("Buy", func() {
		var symbol string = page.GetFormItemByLabel("Stock").(*tview.InputField).GetText()
		cost, costErr := strconv.ParseFloat(costStr, 64)
		quantity, quantityErr := strconv.Atoi(page.GetFormItemByLabel("Quantity").(*tview.InputField).GetText())
		price, priceErr := strconv.ParseFloat(page.GetFormItemByLabel("Price").(*tview.InputField).GetText(), 64)

		var selectedCurrencyId int = getCurrencyIdFromString(selectedCurrency)
		var selectedAccountId int = getAccountIdFromString(selectedAccount)

		if costErr == nil && quantityErr == nil && symbol != "" && priceErr == nil {
			db := connectDb()

			var account Account
			db.First(&account, selectedAccountId)

			var newBalance float64
			if selectedCurrency == "cad" {
				newBalance = account.CadBalance - cost

				if newBalance >= 0 {
					account.CadBalance = newBalance
					updateAppCashCadBalances(appPrimitives, newBalance)
				}

			} else if selectedCurrency == "usd" {
				newBalance = account.UsdBalance - cost

				if newBalance >= 0 {
					account.UsdBalance = newBalance
					updateAppCashUsdBalances(appPrimitives, newBalance)
				}
			}

			// Save the updated record
			if newBalance >= 0 {
				db.Create(&Trade{Symbol: symbol, Quantity: quantity, Price: price, TradeType: "buy", CurrencyID: selectedCurrencyId, AccountID: selectedAccountId})
				db.Save(&account)

				// Query for same stock in specific account to calculate average
				var stock Stock
				result := db.First(&stock, "symbol = ? AND account_id = ?", symbol, selectedAccountId)
				// Update stock record or add new if none found
				if result.RowsAffected != 0 {
					var oldPurchaseValue float64 = calculateStockCost(stock.Quantity, stock.Average)
					var newPurchaseValue float64 = calculateStockCost(quantity, price)
					var newQuantity int = stock.Quantity + quantity
					var newAverage float64 = (oldPurchaseValue + newPurchaseValue) / float64(newQuantity)
					// Update stock record in db
					stock.Average = newAverage
					stock.Quantity = newQuantity
					db.Save(&stock)

				} else {
					db.Create(&Stock{Symbol: symbol, Average: price, Quantity: quantity, CurrencyID: selectedCurrencyId, AccountID: selectedAccountId})
				}

				updatePortfolioStockTable(queryStocks(db, selectedAccountId), appPrimitives)
				appPrimitives.PortfolioDropdown.SetCurrentOption(selectedAccountId)
				pages.SwitchToPage(portfolioPageName)
			}

		}
	})
	page.AddButton("Cancel", func() {
		pages.SwitchToPage(portfolioPageName)
	})

	return page
}

func createPortfolioSellPage(pages *tview.Pages, app *tview.Application, appPrimitives AppPrimitives) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle("Sell")

	var selectedCurrency string
	var selectedAccount string

	page.AddInputField("Stock", "", 20, nil, nil)
	page.AddInputField("Quantity", "", 20, nil, nil)
	page.AddInputField("Price", "", 20, nil, nil)
	page.AddDropDown("Account", []string{"Margin", "TFSA", "RRSP"}, 0, func(option string, optionIndex int) { selectedAccount = option })
	page.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })

	page.AddButton("Sell", func() {

		var symbol string = page.GetFormItemByLabel("Stock").(*tview.InputField).GetText()
		quantity, quantityErr := strconv.Atoi(page.GetFormItemByLabel("Quantity").(*tview.InputField).GetText())
		price, priceErr := strconv.ParseFloat(page.GetFormItemByLabel("Price").(*tview.InputField).GetText(), 64)

		var selectedCurrencyId int = getCurrencyIdFromString(selectedCurrency)
		var selectedAccountId int = getAccountIdFromString(selectedAccount)

		if quantityErr == nil && symbol != "" && priceErr == nil {
			db := connectDb()

			var account Account
			db.First(&account, selectedAccountId)

			// Query stocks for ticker
			stock := queryStock(db, symbol, selectedAccountId)

			if stock.Symbol != "" {
				var newQuantity int = stock.Quantity - quantity
				var cost float64 = calculateStockCost(quantity, price)
				var newBalance float64
				// sell some stock
				if newQuantity > 0 {
					stock.Quantity = newQuantity
					db.Create(&Trade{Symbol: symbol, Quantity: quantity, Price: price, TradeType: "sell", CurrencyID: selectedCurrencyId, AccountID: selectedAccountId})

					if stock.CurrencyID == 1 {
						newBalance = account.CadBalance + cost
						account.CadBalance = newBalance
						updateAppCashCadBalances(appPrimitives, newBalance)
					} else {
						newBalance = account.UsdBalance + cost
						account.UsdBalance = newBalance
						updateAppCashUsdBalances(appPrimitives, newBalance)
					}

					db.Save(&account)
					db.Save(&stock)
					updatePortfolioStockTable(queryStocks(db, selectedAccountId), appPrimitives)

				} else if newQuantity == 0 {
					db.Delete(&stock)
					db.Create(&Trade{Symbol: symbol, Quantity: quantity, Price: price, TradeType: "sell", CurrencyID: selectedCurrencyId, AccountID: selectedAccountId})

					if stock.CurrencyID == 1 {
						newBalance = account.CadBalance + cost
						account.CadBalance = newBalance
						updateAppCashCadBalances(appPrimitives, newBalance)
					} else {
						newBalance = account.UsdBalance + cost
						account.UsdBalance = newBalance
						updateAppCashUsdBalances(appPrimitives, newBalance)
					}

					db.Save(&account)
					db.Save(&stock)
					updatePortfolioStockTable(queryStocks(db, selectedAccountId), appPrimitives)

				}
			}

			pages.SwitchToPage(portfolioPageName)

		}

	})
	page.AddButton("Cancel", func() {
		pages.SwitchToPage(portfolioPageName)
	})

	return page
}

func createFundsPage(pages *tview.Pages) *tview.List {
	fundsPage := tview.NewList()
	fundsPage.SetBorder(true)
	fundsPage.SetTitle("Funds")

	fundsPage.AddItem("Margin", "Some explanatory text", 'a', func() {
		pages.SwitchToPage(fundsMarginPageName)
	})
	fundsPage.AddItem("TFSA", "Some explanatory text", 'b', func() {
		pages.SwitchToPage(fundsTfsaPageName)
	})
	fundsPage.AddItem("RRSP", "Some explanatory text", 'c', func() {
		pages.SwitchToPage(fundsRrspPageName)
	})
	fundsPage.AddItem("Go back", "Press to go back", 'q', func() {
		pages.SwitchToPage(homePageName)
	})

	return fundsPage
}

func createFundsTransactionTypePage(pages *tview.Pages, appPrimitives AppPrimitives, title string) *tview.List {
	page := tview.NewList()
	page.SetBorder(true)
	page.SetTitle(fmt.Sprintf("Funds %s", title))

	db := connectDb()

	if title == "Margin" {
		page.AddItem("Deposit", "Add funds to your margin account", 'a', func() { pages.SwitchToPage(fundsMarginDepositPageName) })
		page.AddItem("Withdraw", "Remove funds to your margin account", 'b', func() {

			var account Account
			db.First(&account, 1)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage(fundsMarginWithdrawPageName)
		})

	} else if title == "TFSA" {
		page.AddItem("Deposit", "Add funds to your TFSA account", 'a', func() { pages.SwitchToPage(fundsTfsaDepositPageName) })
		page.AddItem("Withdraw", "Remove funds to your TFSA account", 'b', func() {

			var account Account
			db.First(&account, 2)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage(fundsTfsaWithdrawPageName)
		})

	} else if title == "RRSP" {
		page.AddItem("Deposit", "Add funds to your RRSP account", 'a', func() { pages.SwitchToPage(fundsRrspDepositPageName) })
		page.AddItem("Withdraw", "Remove funds to your RRSP account", 'b', func() {

			var account Account
			db.First(&account, 3)

			updateAppCashCadBalances(appPrimitives, account.CadBalance)
			updateAppCashUsdBalances(appPrimitives, account.UsdBalance)
			pages.SwitchToPage(fundsRrspWithdrawPageName)
		})
	}

	page.AddItem("Go back", "Press to go back", 'q', func() { pages.SwitchToPage(fundsPageName) })
	return page
}

func createFundsDepositPage(pages *tview.Pages, appPrimitives AppPrimitives, title string) *tview.Form {
	page := tview.NewForm()
	page.SetBorder(true)
	page.SetTitle(fmt.Sprintf("Deposit Funds %s", title))

	var selectedCurrency string = ""

	page.AddInputField("Amount", "", 20, nil, nil)
	page.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })

	page.AddButton("Save", func() {
		db := connectDb()

		amountStr := page.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amount, err := strconv.ParseFloat(amountStr, 64)

		// Set the proper account ID based on the title
		var accountId int = getAccountIdFromString(title)

		// Create transaction
		if err == nil {
			db.Create(&CashTransaction{Amount: amount, Type: "deposit", CurrencyID: getCurrencyIdFromString(selectedCurrency), AccountID: accountId})

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
		switchToFundsAccountOptions(title, pages)

	})

	page.AddButton("Cancel", func() {
		page.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")
		switchToFundsAccountOptions(title, pages)
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

	db := connectDb()

	balanceFlex.AddItem(appPrimitives.FundsCadTextView, 0, 1, false)
	balanceFlex.AddItem(appPrimitives.FundsUsdTextView, 0, 1, false)

	var selectedCurrency string = ""

	form := tview.NewForm()
	form.SetBorder(true)

	form.AddInputField("Amount", "", 20, nil, nil)
	form.AddDropDown("Currency", []string{"cad", "usd"}, 0, func(option string, optionIndex int) { selectedCurrency = option })

	form.AddButton("Save", func() {

		amountStr := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amount, err := strconv.ParseFloat(amountStr, 64)

		// Set the proper account ID based on the title
		var accountId int = getAccountIdFromString(title)

		// Create transaction
		if err == nil {

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
				db.Create(&CashTransaction{Amount: amount, Type: "withdraw", CurrencyID: getCurrencyIdFromString(selectedCurrency), AccountID: accountId})
				db.Save(&account)
			}

		}

		form.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")
		switchToFundsAccountOptions(title, pages)

	})

	form.AddButton("Cancel", func() {
		form.GetFormItemByLabel("Amount").(*tview.InputField).SetText("")
		switchToFundsAccountOptions(title, pages)

	})

	page.AddItem(balanceFlex, 0, 1, false)
	page.AddItem(form, 0, 4, false)

	// Handle key presses -= Navigation
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(form)

		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			switchToFundsAccountOptions(title, pages)
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
