package sldemo

import (
	"log/slog"

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

func (acc Account) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", acc.Id.String()),
		slog.String("name", "[REDACTED]"),
		slog.String("bank_code", acc.BankCode),
		slog.String("branch_code", acc.BranchCode),
		slog.String("account_number", "[REDACTED]"),
		slog.String("account_type", acc.AccountType),
	)
}

func (acc AccountRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", "[REDACTED]"),
		slog.String("bank_code", acc.BankCode),
		slog.String("branch_code", acc.BranchCode),
		slog.String("account_number", "[REDACTED]"),
		slog.String("account_type", acc.AccountType),
	)
}
