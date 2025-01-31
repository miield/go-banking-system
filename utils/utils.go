package utils

import (
	"fmt"
	"github.com/miield/go-banking-system/accounts"
	"go-banking-system/accounts"
	"github.com/xuri/excelize/v2"
)

// all account
var allAccount = []accounts.Account{}

func DisplayAllAccountsToExcel() error {
	f := excelize.NewFile()

	// Create header
	f.SetCellValue("Sheet1", "A1", "Account Number")
	f.SetCellValue("Sheet1", "B1", "Account Name")
	f.SetCellValue("Sheet1", "C1", "Balance")

	for i, acc := range allAccount {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), acc.AccountNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), acc.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), acc.Balance)
	}

	// Save to file
	if err := f.SaveAs("accounts.xlsx"); err != nil {
		return fmt.Errorf("failed to save excel file: %s", err)
	}

	fmt.Println("Accounts data saved to accounts.xlsx successfully!")
	return nil
}