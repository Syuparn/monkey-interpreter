package runscript

import (
	"fmt"
	"testing"
)

func TestRunScriptErrors(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		// reading error
		{
			"sample/notexistingfile.monkey",
			"file could not open: sample/notexistingfile.monkey",
		},
		// evaluator error
		{
			"sample/err_no_ident.monkey",
			"identifier not found: a",
		},
		// parser error
		{
			"sample/err_fail2parse.monkey",
			(MONKEY_FACE +
				"Woops! We ran into some monkey business here!\n" +
				" parser errors:\n" +
				"\tno prefix parse function for + found\n"),
		},
	}

	for _, tt := range tests {
		_, err := RunScript(tt.fileName)

		if err == nil {
			t.Fatalf("err should occur. fileName=%s", tt.fileName)
			return
		}

		if fmt.Sprintf("%s", err) != tt.expected {
			t.Fatalf("error message was wrong. got=\n%s\n, expected=\n%s\n",
				fmt.Sprintf("%s", err), tt.expected)
		}
	}
}
