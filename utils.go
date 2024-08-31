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

func updateAppMarketValue(appPrimitives AppPrimitives, marketValue float64) {
	appPrimitives.PortfolioMarketValue.SetText(fmt.Sprintf("Market Value: %.2f", marketValue))
}

func updateAppRealizedPLValue(appPrimitives AppPrimitives, value float64) {
	appPrimitives.RealizedPLValue.SetText(fmt.Sprintf("Realized P/L: %.2f", value))
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

	stockMap := make(map[string]StockInfo)

	for _, stock := range stocks {
		if item, ok := stockMap[stock.Symbol]; ok {
			// Symbol exists; update the existing StockInfo
			item.Value = (stock.Average * float64(stock.Quantity)) + item.Value
			item.Quantity = stock.Quantity + item.Quantity
			item.Average = ((stock.Average * float64(stock.Quantity)) + item.Value) / (float64(stock.Quantity + item.Quantity))
			stockMap[stock.Symbol] = item
		} else {
			// Symbol does not exist; create a new StockInfo
			stockMap[stock.Symbol] = StockInfo{Symbol: stock.Symbol, Value: stock.Average * float64(stock.Quantity), Quantity: stock.Quantity, Average: stock.Average}
		}

	}

	var row int = 1
	for key := range stockMap {
		appPrimitives.PortfolioStockTable.SetCell(row, 0, tview.NewTableCell(key))
		appPrimitives.PortfolioStockTable.SetCell(row, 1, tview.NewTableCell(strconv.Itoa(stockMap[key].Quantity)))
		appPrimitives.PortfolioStockTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.2f", stockMap[key].Average)))
		row++
	}
}

func calculateStockCost(quantity int, price float64) float64 {
	return float64(quantity) * price
}

func calculateRealizedPL(sellPrice float64, average float64, quantity int) float64 {
	return (sellPrice - average) * float64(quantity)
}
