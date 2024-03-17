package account

import (
	"github.com/ippoippo/slog-lt/zapdemo"
	"github.com/ippoippo/slog-lt/zapdemo/persistence"
)

func NewStorage() *persistence.Storage[*zapdemo.Account] {
	return persistence.New[*zapdemo.Account]()
}
