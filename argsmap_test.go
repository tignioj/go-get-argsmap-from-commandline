package argsmap

import (
	"fmt"
	"testing"
)

func TestGetCommandLineArgMap(t *testing.T) {
	args := []string{"filename", "-flag1", "user_input_v1", "-flag2"}
	clobj, err := NewCommandLineObj("test.json", args)
	if err != nil {
		t.Fatal(args, err)
	} else {
		clobj.ShowHelp()
		fmt.Println(clobj.GetCommandLineMap)
	}

}
