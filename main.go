package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
)

var modules []module.Module

func main() {
	keymap := make(map[string]string)
	keymap["pascual"] = "victor"
	keymap["nombre"] = "pedro"

	input := "El @{nombre} del pascual es @{pascual}"

	out, err := module.RenderString(input, keymap)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(out)
	}

	init := jmdict.InitOptions{}
	init.DictionaryPath = "./Repository/Vocabulary/JMdict_e_examp.xml"
	init.FormatterPath = "./Tools/POLXML/out.xml"

	mod, err := jmdict.Initialize(init)

	if err != nil {
		fmt.Println(err)
		return
	}

	mod.Demo()
}
