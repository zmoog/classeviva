package feedback

import (
	"fmt"
	"io"
	"os"
)

type Feedback struct {
	out io.Writer
	err io.Writer
}

func New(out, err io.Writer) *Feedback {
	return &Feedback{out: out, err: err}
}

func Default() *Feedback {
	return New(os.Stdout, os.Stderr)
}

func (fb *Feedback) Println(v interface{}) {
	fmt.Fprintln(fb.out, v)
}

// type Feedback interface {
// 	Println(a ...interface{}) (n int, err error)
// }

// type Console struct {
// }

// func (c Console) Println(a ...interface{}) (n int, err error) {
// 	return fmt.Println(a...)
// }
