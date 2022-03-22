package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
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
}
