package feedback

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	Text OutputFormat = iota
	JSON
)

type OutputFormat int

type Feedback struct {
	out    io.Writer
	err    io.Writer
	format OutputFormat
}

func New(out, err io.Writer, format OutputFormat) *Feedback {
	return &Feedback{out: out, err: err, format: format}
}

func Default() *Feedback {
	return New(os.Stdout, os.Stderr, Text)
}

func (fb *Feedback) SetFormat(format OutputFormat) {
	fb.format = format
}

func (fb *Feedback) Println(v interface{}) {
	fmt.Fprintln(fb.out, v)
}

func (fb *Feedback) Error(v interface{}) {
	fmt.Fprintln(fb.err, v)
}

func (fb *Feedback) PrintResult(result Result) error {
	var output string

	switch fb.format {
	case JSON:
		byteOutput, err := json.MarshalIndent(result.Data(), "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshall result: %w", err)
		}
		output = string(byteOutput)
	default:
		output = result.String()
	}

	_, err := fmt.Fprint(fb.out, output)
	if err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	return nil
}

type Result interface {
	fmt.Stringer
	Data() interface{}
}
