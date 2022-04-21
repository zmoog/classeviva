package commands

import (
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type Command interface {
	Execute(adapter spaggiari.Adapter) error
}
