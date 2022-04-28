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

func (fb *Feedback) PrintResult(result Result) (err error) {
	switch fb.format {
	case JSON:
		output, _ := json.MarshalIndent(result.Data(), "", "  ")
		_, err = fmt.Fprint(fb.out, string(output))
	default:
		_, err = fmt.Fprint(fb.out, result.String())
	}
	return err
}

type Result interface {
	fmt.Stringer
	Data() interface{}
}
