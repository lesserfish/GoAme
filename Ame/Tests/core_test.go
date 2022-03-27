package ametest

import(
	"testing"
	"github.com/lesserfish/GoAme/Modules/Core"
)

func TestA(t *testing.T){
	core.Demo()
	t.Errorf("Bullshit error!")
}

