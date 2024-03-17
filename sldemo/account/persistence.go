package account

import (
	"github.com/ippoippo/slog-lt/sldemo"
	"github.com/ippoippo/slog-lt/sldemo/persistence"
)

func NewStorage() *persistence.Storage[*sldemo.Account] {
	return persistence.New[*sldemo.Account]()
}
