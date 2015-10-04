package kovacs

import "testing"

func TestStdinTypes(t *testing.T) {
	// just make sure these compile
	var _ StdinType = StdinDevNull
	var _ StdinType = StdinNamePerLine
	var _ StdinType = StdinArray{"sdfkj", "sdflkjsdf"}
}
