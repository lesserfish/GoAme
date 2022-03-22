package jmdict

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
)

type Parser struct {
	DictionaryPath string
	dictionary     JMdict
}
type InitOptions struct {
	DictionaryPath string
}

func Initialize(options InitOptions) (module.Module, error) {
	newParser := new(Parser)
	newParser.DictionaryPath = options.DictionaryPath

	err := LoadDictionary(newParser)

	return *newParser, err
}
func (parser Parser) Demo() {
	fmt.Println("Hello from JMdict from ", parser.DictionaryPath)
}
