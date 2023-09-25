package lexfn

import "testing"

type binTest struct {
	Name           string
	Input          string
	ExpectedOutput bool
}

type binTests []binTest

var (
	isLiteralInputTests = binTests{
		{"Blank", "", false},
		{"Single Dot", ".", false},
	}
	isRDFLiteralTests = append(
		isLiteralInputTests,
		binTest{"Random text string", `"random string"`, true},
	)
	isNumericLiteralTests = append(
		isLiteralInputTests,
		binTest{"Random text string", `"random string"`, false},
	)
)

func TestIsLiteral(t *testing.T) {
	for _, tc := range isLiteralInputTests {
		if isLiteral(tc.Input) != tc.ExpectedOutput {
			t.Errorf("%v test fail", tc.Name)
		}
	}
}

func TestRDFLiteral(t *testing.T) {
	for _, tc := range isRDFLiteralTests {
		if isRDFLiteral(tc.Input) != tc.ExpectedOutput {
			t.Errorf("%v test fail", tc.Name)
		}
	}
}

func TestNumericLiteral(t *testing.T) {
	for _, tc := range isNumericLiteralTests {
		if isNumericLiteral(tc.Input) != tc.ExpectedOutput {
			t.Errorf("%v test fail", tc.Name)
		}
	}
}

func TestBooleanliteral(t *testing.T) {
	if isBooleanLiteral(".") {
		t.Error(". returning true")
	}
}

func TestInteger(t *testing.T) {
	if isInteger(".") {
		t.Error(". returning true")
	}
}
func TestDecimal(t *testing.T) {
	if isDecimal(".") {
		t.Error(". returning true")
	}
}
