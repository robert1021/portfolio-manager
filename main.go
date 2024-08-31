package main

import (
	"os"

	"github.com/rivo/tview"
	"gorm.io/gorm"
)

// Account represents an account with different currencies
type Account struct {
	gorm.Model
	AccountType string `gorm:"not null"`
	CadBalance  float64
	UsdBalance  float64
}

type Currency struct {
	gorm.Model
	CurrenyType string `gorm:"not null"`
}

// Stock represents a stock with an associated account
type Stock struct {
	gorm.Model
	Symbol     string   `gorm:"not null"`
	Average    float64  `gorm:"not null"`
	Quantity   int      `gorm:"not null"`
	CurrencyID int      // Foreign key to Currency
	AccountID  int      // Foreign key to Account
	Account    Account  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Currency   Currency `gorm:"foreignKey:CurrencyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// CashTransaction represents a cash transaction linked to an account
type CashTransaction struct {
	gorm.Model
	Amount     float64  `gorm:"not null"`
	Type       string   `gorm:"not null"`
	CurrencyID int      // Foreign key to Currency
	AccountID  int      // Foreign key to Account
	Account    Account  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Currency   Currency `gorm:"foreignKey:CurrencyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Trade represents a trade transaction linked to an account
type Trade struct {
	gorm.Model
	Symbol     string   `gorm:"not null"`
	Quantity   int      `gorm:"not null"`
	Price      float64  `gorm:"not null"`
	TradeType  string   `gorm:"not null"` // Type of trade: "buy" or "sell"
	RealizedPL float64  //store realized P&L
	CurrencyID int      // Foreign key to Currency
	AccountID  int      // Foreign key to Account
	Account    Account  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Currency   Currency `gorm:"foreignKey:CurrencyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type StockInfo struct {
	Symbol   string
	Value    float64
	Quantity int
	Average  float64
}

type AppPrimitives struct {
	FundsCadTextView     *tview.TextView
	FundsUsdTextView     *tview.TextView
	PortfolioCadTextView *tview.TextView
	PortfolioUsdTextView *tview.TextView
	BuyStockModal        *tview.Modal
	PortfolioStockTable  *tview.Table
	PortfolioDropdown    *tview.DropDown
	PortfolioMarketValue *tview.TextView
	RealizedPLValue      *tview.TextView
}

func main() {

	// Set up db
	if _, err := os.Stat(dbFilename); err != nil {
		db := connectDb()

		db.AutoMigrate(&Account{})
		db.AutoMigrate(&Stock{})
		db.AutoMigrate(&CashTransaction{})
		db.AutoMigrate(&Trade{})

		// Set up accounts
		db.Create(&Account{AccountType: "margin", CadBalance: 0, UsdBalance: 0})
		db.Create(&Account{AccountType: "tfsa", CadBalance: 0, UsdBalance: 0})
		db.Create(&Account{AccountType: "rrsp", CadBalance: 0, UsdBalance: 0})

		// Set up currencies
		db.Create(&Currency{CurrenyType: "cad"})
		db.Create(&Currency{CurrenyType: "usd"})

	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	// Hold primitives that can be passed around the app
	appPrimitives := AppPrimitives{
		FundsCadTextView:     tview.NewTextView(),
		FundsUsdTextView:     tview.NewTextView(),
		PortfolioCadTextView: tview.NewTextView(),
		PortfolioUsdTextView: tview.NewTextView(),
		BuyStockModal:        tview.NewModal(),
		PortfolioStockTable:  tview.NewTable(),
		PortfolioDropdown:    tview.NewDropDown(),
		PortfolioMarketValue: tview.NewTextView(),
		RealizedPLValue:      tview.NewTextView(),
	}

	// Home page
	homePage := createHomePage(pages, func() { app.Stop() }, appPrimitives)
	// Portfolio Page
	portfolioPage := createPortfolioPage(pages, app, appPrimitives)
	portfolioBuyPage := createPortfolioBuyPage(pages, app, appPrimitives)
	portfolioSellPage := createPortfolioSellPage(pages, app, appPrimitives)
	// Add Funds Page
	fundsPage := createFundsPage(pages)
	fundsMarginPage := createFundsTransactionTypePage(pages, appPrimitives, "Margin")
	fundsMarginDepositPage := createFundsDepositPage(pages, appPrimitives, "Margin")
	fundsMarginWithdrawPage := createFundsWithdrawPage(pages, app, appPrimitives, "Margin")
	fundsTfsaPage := createFundsTransactionTypePage(pages, appPrimitives, "TFSA")
	fundsTfsaDepositPage := createFundsDepositPage(pages, appPrimitives, "TFSA")
	fundsTfsaWithdrawPage := createFundsWithdrawPage(pages, app, appPrimitives, "TFSA")
	fundsRrspPage := createFundsTransactionTypePage(pages, appPrimitives, "RRSP")
	fundsRrspDepositPage := createFundsDepositPage(pages, appPrimitives, "RRSP")
	fundsRrspWithdrawPage := createFundsWithdrawPage(pages, app, appPrimitives, "RRSP")
	// Statisitcs Page
	statisticsPage := createStatisticsPage(pages)

	// Add pages to Pages
	pages.AddPage(homePageName, homePage, true, true)
	pages.AddPage(portfolioPageName, portfolioPage, true, false)
	pages.AddPage(portfolioBuyPageName, portfolioBuyPage, true, false)
	pages.AddPage(portfolioSellPageName, portfolioSellPage, true, false)
	pages.AddPage(fundsPageName, fundsPage, true, false)
	pages.AddPage(fundsMarginPageName, fundsMarginPage, true, false)
	pages.AddPage(fundsMarginDepositPageName, fundsMarginDepositPage, true, false)
	pages.AddPage(fundsMarginWithdrawPageName, fundsMarginWithdrawPage, true, false)
	pages.AddPage(fundsTfsaPageName, fundsTfsaPage, true, false)
	pages.AddPage(fundsTfsaDepositPageName, fundsTfsaDepositPage, true, false)
	pages.AddPage(fundsTfsaWithdrawPageName, fundsTfsaWithdrawPage, true, false)
	pages.AddPage(fundsRrspPageName, fundsRrspPage, true, false)
	pages.AddPage(fundsRrspDepositPageName, fundsRrspDepositPage, true, false)
	pages.AddPage(fundsRrspWithdrawPageName, fundsRrspWithdrawPage, true, false)
	pages.AddPage(statisticsPageName, statisticsPage, true, false)

	// Set the initial page to the home page
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}
