package cnamount

import (
	"testing"
	"fmt"
)

func testAmount(f float64) {
	fmt.Printf("%.2f => %s\n", f, ToCNAmount(f))
}

func TestCNAmount(t *testing.T) {
	testAmount(435235324532.0)
	testAmount(100.02)
	testAmount(10100.02)
	testAmount(340210100.02)
	testAmount(3400000000.02)
	testAmount(4352352343400000000.02)
	testAmount(.00)
	testAmount(.12)
	testAmount(9999.0)
	testAmount(19800.0)
	testAmount(2980.0)
	testAmount(500200.0)
	testAmount(103.0)
}

