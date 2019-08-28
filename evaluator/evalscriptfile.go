package evaluator

import (
	"../lexer"
	"../object"
	"../parser"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func EvalScriptFile(fileName string) (*object.Environment, error) {
	// NOTE: enviroment内の変数"THIS_DIR","THIS_FILE"に
	// スクリプトファイルのディレクトリ/ファイルパスをSTRINGで格納

	script, absFileName, err := tryAllPathsReadScript(fileName)
	if err != nil {
		return nil, err
	}

	l := lexer.New(script)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return nil, formatParserErrors(p.Errors())
	}

	env := object.NewEnvironment()
	// env内にディレクトリ場所を格納する変数"THIS_DIR"を束縛
	if absFileName != "" {
		env.Set("THIS_DIR", &object.String{Value: filepath.Dir(absFileName) + "/"})
		env.Set("THIS_FILE", &object.String{Value: absFileName})
	}

	evaluated := Eval(program, env)

	if errObj, ok := evaluated.(*object.Error); ok {
		return nil, formatEvaluatorErrors(errObj)
	}

	return env, nil
}

func tryAllPathsReadScript(fileName string) (string, string, error) {
	script, err := readScript(fileName)
	// if found, return it
	if err == nil {
		// make separators ("/" or "\") all the same
		return script, filepath.Clean(fileName), nil
	}

	candidatePaths, pathErr := defaultPaths()
	if pathErr != nil {
		return "", "", fmt.Errorf("fail to create defaultPaths:\n%s", err)
	}

	// try again!
	for _, path := range candidatePaths {
		absFileName := filepath.Join(path, fileName) // Join cleans path inside
		script, err := readScript(absFileName)
		if err == nil {
			return script, absFileName, nil
		}
	}

	return "", "", err
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

var defaultPathCache = []string{}

func defaultPaths() ([]string, error) {
	if len(defaultPathCache) != 0 {
		return defaultPathCache, nil
	}

	relativePaths := []string{
		"scripts",
	}

	monkeyPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	mainDir := filepath.Dir(monkeyPath)

	paths := []string{}
	for _, relativePath := range relativePaths {
		paths = append(paths, filepath.Join(mainDir, relativePath))
	}

	if len(defaultPathCache) == 0 {
		defaultPathCache = paths
	}

	return paths, nil
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
