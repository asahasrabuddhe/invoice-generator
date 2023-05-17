package chrome

import (
	"os/exec"
	"testing"
)

func TestLocate(t *testing.T) {
	if exe := Executable(); exe == "" {
		t.Fatal()
	} else {
		t.Log(exe)
		b, err := exec.Command(exe, "--version").CombinedOutput()
		t.Log(string(b))
		t.Log(err)
	}
}
