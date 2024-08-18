package main

import (
	"os"

	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
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
	Postion    int      `gorm:"not null"`
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
	Quantity   float64  `gorm:"not null"`
	Price      float64  `gorm:"not null"`
	TradeType  string   `gorm:"not null"` // Type of trade: "buy" or "sell"
	CurrencyID int      // Foreign key to Currency
	AccountID  int      // Foreign key to Account
	Account    Account  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Currency   Currency `gorm:"foreignKey:CurrencyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func main() {

	// Set up db
	if _, err := os.Stat("portfolio-manager.db"); err != nil {
		db, err := gorm.Open(sqlite.Open("portfolio-manager.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

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

	// Home page
	homePage := createHomePage(pages, func() { app.Stop() })
	// Portfolio Page
	portfolioPage := createPortfolioPage(pages, app)
	// Add Funds Page
	fundsPage := createFundsPage(pages)
	fundsMarginPage := createFundsMarginPage(pages)
	fundsMarginDepositPage := createFundsMarginDepositPage(pages)
	fundsTfsaPage := createFundsTfsaPage(pages)
	fundsRrspPage := createFundsRrspPage(pages)
	// Statisitcs Page
	statisticsPage := createStatisticsPage(pages)

	// Add pages to Pages
	pages.AddPage("home", homePage, true, true)
	pages.AddPage("portfolio", portfolioPage, true, false)
	pages.AddPage("funds", fundsPage, true, false)
	pages.AddPage("fundsMargin", fundsMarginPage, true, false)
	pages.AddPage("fundsMarginDeposit", fundsMarginDepositPage, true, false)
	pages.AddPage("fundsTfsa", fundsTfsaPage, true, false)
	pages.AddPage("fundsRrsp", fundsRrspPage, true, false)
	pages.AddPage("statistics", statisticsPage, true, false)

	// Set the initial page to the home page
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

}
