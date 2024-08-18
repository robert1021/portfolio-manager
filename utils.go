package main

import "fmt"

func updateAppCashCadBalances(appPrimitives AppPrimitives, balance float64) {
	appPrimitives.FundsCadTextView.SetText(fmt.Sprintf("Balance CAD: %.2f", balance))
	appPrimitives.PortfolioCadTextView.SetText(fmt.Sprintf("Balance CAD: %.2f", balance))
}

func updateAppCashUsdBalances(appPrimitives AppPrimitives, balance float64) {
	appPrimitives.FundsUsdTextView.SetText(fmt.Sprintf("Balance USD: %.2f", balance))
	appPrimitives.PortfolioUsdTextView.SetText(fmt.Sprintf("Balance USD: %.2f", balance))
}
