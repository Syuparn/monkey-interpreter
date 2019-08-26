package evaluator

import (
	"../lexer"
	"../object"
	"../parser"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func EvalScriptFile(fileName string) (*object.Environment, error) {
	const DEFAULT_PATH = "../scripts/"

	scriptCode, err := readScript(fileName)

	if err != nil {
		// 読み込み失敗した場合は、デフォルトパスをprefixに付けてもう一度探す
		alternativeCode, alternativeErr := readScript(DEFAULT_PATH + fileName)
		if alternativeErr != nil {
			return nil, err
		} else {
			// 代わりのファイル名で読み込めた場合はそちらのコードを使用
			scriptCode = alternativeCode
		}
	}

	l := lexer.New(scriptCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return nil, formatParserErrors(p.Errors())
	}

	env := object.NewEnvironment()
	evaluated := Eval(program, env)

	if errObj, ok := evaluated.(*object.Error); ok {
		return nil, formatEvaluatorErrors(errObj)
	}

	return env, nil
}

func readScript(fileName string) (string, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return "", fmt.Errorf("file could not open: %s", fileName)
	}

	script, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("file could not read: %s", fileName)
	}

	return string(script), nil
}

const MONKEY_FACE = `
           __,__
  .--.  .-"     "-.  .--.
 / .. \/  .-. .-.  \/ .. \
| |  '|  /   Y   \  |'  | |
| \   \  \ 0 | 0 /  /   / |
 \ '- ,\.-"""""""-./, -' /
  ''-' /_   ^ ^   _\ '-''
      |  \._   _./  |
      \   \ '~' /   /
       '._ '-=-' _.'
          '-----'
`

func formatParserErrors(errors []string) error {
	var errMsg bytes.Buffer
	errMsg.WriteString(MONKEY_FACE)
	errMsg.WriteString("Woops! We ran into some monkey business here!\n")
	errMsg.WriteString(" parser errors:\n")
	for _, err := range errors {
		errMsg.WriteString("\t" + err + "\n")
	}

	return fmt.Errorf(errMsg.String())
}

func formatEvaluatorErrors(errObj *object.Error) error {
	return fmt.Errorf(errObj.Message)
}
