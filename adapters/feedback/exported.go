package feedback

var (
	fb = Default()
)

func SetDefault(f *Feedback) {
	fb = f
}

func SetFormat(format OutputFormat) {
	fb.SetFormat(format)
}

func Println(v interface{}) {
	fb.Println(v)
}

func Error(v interface{}) {
	fb.Error(v)
}

func PrintResult(result Result) (err error) {
	return fb.PrintResult(result)
}
