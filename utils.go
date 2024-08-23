package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

func updateAppCashCadBalances(appPrimitives AppPrimitives, balance float64) {
	appPrimitives.FundsCadTextView.SetText(fmt.Sprintf("Balance CAD: %.2f", balance))
	appPrimitives.PortfolioCadTextView.SetText(fmt.Sprintf("Balance CAD: %.2f", balance))
}

func updateAppCashUsdBalances(appPrimitives AppPrimitives, balance float64) {
	appPrimitives.FundsUsdTextView.SetText(fmt.Sprintf("Balance USD: %.2f", balance))
	appPrimitives.PortfolioUsdTextView.SetText(fmt.Sprintf("Balance USD: %.2f", balance))
}

func getCurrencyIdFromString(selectedCurrency string) int {
	var currencyId int
	if selectedCurrency == "cad" {
		currencyId = 1
	} else if selectedCurrency == "usd" {
		currencyId = 2
	}

	return currencyId
}

func getAccountIdFromString(accountType string) int {
	var stringLower string = strings.ToLower(accountType)
	var accountId int

	switch stringLower {
	case "margin":
		accountId = 1
	case "tfsa":
		accountId = 2
	case "rrsp":
		accountId = 3
	}

	return accountId
}

func switchToFundsAccountOptions(accountType string, pages *tview.Pages) {
	var stringLower string = strings.ToLower(accountType)

	switch stringLower {
	case "margin":
		pages.SwitchToPage(fundsMarginPageName)
	case "tfsa":
		pages.SwitchToPage(fundsTfsaPageName)
	case "rrsp":
		pages.SwitchToPage(fundsRrspPageName)
	}
}

func updatePortfolioStockTable(stocks []Stock, appPrimitives AppPrimitives) {
	appPrimitives.PortfolioStockTable.Clear()
	// Set up header columns
	appPrimitives.PortfolioStockTable.SetCell(0, 0, tview.NewTableCell("STOCK"))
	appPrimitives.PortfolioStockTable.SetCell(0, 1, tview.NewTableCell("QUANTITY"))
	appPrimitives.PortfolioStockTable.SetCell(0, 2, tview.NewTableCell("AVERAGE"))

	var row int = 1
	for _, stock := range stocks {
		appPrimitives.PortfolioStockTable.SetCell(row, 0, tview.NewTableCell(stock.Symbol))
		appPrimitives.PortfolioStockTable.SetCell(row, 1, tview.NewTableCell(strconv.Itoa(stock.Quantity)))
		appPrimitives.PortfolioStockTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.2f", stock.Average)))
		row++
	}
}
