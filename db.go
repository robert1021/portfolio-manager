package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func queryStock(db *gorm.DB, symbol string, accountId int) Stock {
	var stock Stock
	db.First(&stock, "symbol = ? AND account_id = ?", symbol, accountId)
	return stock
}

func queryStocks(db *gorm.DB, accountId int) []Stock {
	var stocks []Stock
	db.Find(&stocks, "account_id = ?", accountId)
	return stocks
}

func queryAllStocks(db *gorm.DB) []Stock {
	var stocks []Stock
	db.Find(&stocks)
	return stocks
}

func queryAccountMarketValue(db *gorm.DB, accountId int) float64 {
	var stocks []Stock = queryStocks(db, accountId)
	var marketValue float64

	for _, stock := range stocks {
		marketValue += calculateStockCost(stock.Quantity, stock.Average)
	}
	return marketValue
}

func queryAllAccountsMarketValue(db *gorm.DB) float64 {
	return queryAccountMarketValue(db, 1) + queryAccountMarketValue(db, 2) + queryAccountMarketValue(db, 3)
}
