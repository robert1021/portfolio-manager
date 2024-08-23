package main

import (
	"fmt"
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
