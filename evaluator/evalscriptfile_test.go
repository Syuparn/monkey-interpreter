package evaluator

import (
	"fmt"
	"os"
	"testing"
)

func TestRunScriptErrors(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("fail to fetch current diretory:\n%s", err)
	}

	tests := []struct {
		fileName string
		expected string
	}{
		// reading error
		{
			"errsample/notexistingfile.monkey",
			"file could not open: " + curDir + "/errsample/notexistingfile.monkey",
		},
		// parser error
		{
			"errsample/err_fail2parse.monkey",
			(MONKEY_FACE +
				"Woops! We ran into some monkey business here!\n" +
				" parser errors:\n" +
				"\tno prefix parse function for + found\n"),
		},
		// evaluator error
		{
			"errsample/err_no_ident.monkey",
			"identifier not found: a",
		},
	}

	for _, tt := range tests {
		_, err := EvalScriptFile(curDir + "/" + tt.fileName)

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
