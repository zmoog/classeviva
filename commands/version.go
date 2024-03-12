package commands

import (
	"fmt"

	"github.com/zmoog/classeviva/adapters/feedback"
)

var (
	version string = "v0.0.0"
	commit  string = "123"
	date    string = "2022-05-08"
	builtBy string = "zmoog"
)

type VersionCommand struct{}

func (c VersionCommand) ExecuteWith(uow UnitOfWork) error {
	return feedback.PrintResult(VersionResult{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	})
}

type VersionResult struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"built_by"`
}

func (r VersionResult) String() string {
	return fmt.Sprintf("Classeviva CLI %v (%v) %v by %v\n", version, commit, date, builtBy)
}

func (r VersionResult) Data() interface{} {
	return r
}
