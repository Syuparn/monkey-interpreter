package runscript

import (
	"../evaluator"
	"fmt"
	"io"
)

func RunScript(fileName string, out io.Writer) {
	_, err := evaluator.EvalScriptFile(fileName)
	if err != nil {
		io.WriteString(out, fmt.Sprintf("%s", err))
	}
}
