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

	mod, err := jmdict.Initialize(jmdict.InitOptions{"./Repository/Vocabulary/JMdict_e_examp.xml", ""})
	mod.Demo()
	if err != nil {
		fmt.Println(err)
		return
	}

	card := module.Card{}
	kin := []string{"食べる"}

	e := mod.Render(kin, &card)

	fmt.Println(e)
}
