package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	core "github.com/lesserfish/GoAme/Modules/Core"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
)

var modules []module.Module

func main() {
	Configure()
	core.Demo()
	mod, err := jmdict.Initialize(jmdict.InitOptions{"/home/lesserfish/Documents/tmp/godemo/JMdict.xml"})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Loaded dictionary!")
	}

	modules = append(modules, mod)
}
