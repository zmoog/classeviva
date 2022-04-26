package commands

import (
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type Command interface {
	ExecuteWith(uow UnitOfWork) error
}

type UnitOfWork struct {
	Adapter  spaggiari.Adapter
	Feedback *feedback.Feedback
}
