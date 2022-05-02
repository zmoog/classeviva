package commands_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// UnmarshalFrom unmarshals test data stored in JSON file located
// at `path` into the given `data` struct(ure).
func UnmarshalFrom(path string, v interface{}) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("can't unmarshal test data from %v: %w", path, err)
	}

	if !stat.Mode().IsRegular() {
		return fmt.Errorf("can't read from %v", path)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read test data from %v: %w", path, err)
	}

	err = json.Unmarshal(content, v)
	if err != nil {
		return fmt.Errorf("can't unmarshal test data from %v: %w", path, err)
	}

	return nil
}
