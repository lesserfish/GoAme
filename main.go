package main

import (
	"github.com/lesserfish/GoAme/Modules/Core"
	"github.com/lesserfish/GoAme/Modules/JMDict"
)

func main(){
	Configure()
	core.Demo()
	jmdict.Demo()
	jmdict.Call()
}
