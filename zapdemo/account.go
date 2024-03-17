package zapdemo

import (
	"fmt"

	"github.com/google/uuid"
)

type AccountRequest struct {
	Name          string `json:"name"`
	BankCode      string `json:"bank_code"`
	BranchCode    string `json:"branch_code"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
}

type Account struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	BankCode      string    `json:"bank_code"`
	BranchCode    string    `json:"branch_code"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
}

func AccountFromRequest(id uuid.UUID, a AccountRequest) *Account {

	return &Account{
		Id:            id,
		Name:          a.Name,
		BankCode:      a.BankCode,
		BranchCode:    a.BranchCode,
		AccountNumber: a.AccountNumber,
		AccountType:   a.AccountType,
	}
}

// Zap sensistive data hiding
func (acc Account) String() string {
	return fmt.Sprintf(
		"[id:%v - bankCode:%v - branchCode:%v - accountType:%v",
		acc.Id,
		acc.BankCode,
		acc.BranchCode,
		acc.AccountType,
	)
}

func (acc AccountRequest) String() string {
	return fmt.Sprintf(
		"[bankCode:%v - branchCode:%v - accountType:%v",
		acc.BankCode,
		acc.BranchCode,
		acc.AccountType,
	)
}
