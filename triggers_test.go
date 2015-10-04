package kovacs

import "testing"

func TestStdinTypes(t *testing.T) {
	tr := &Trigger{}

	// just make sure these compile
	tr.Stdin = StdinDevNull
	tr.Stdin = StdinNamePerLine
	tr.Stdin = StdinArray{"sdfkj", "sdflkjsdf"}

	// and these do not
	// tr.Stdin = "sdflkjsdf"
}
