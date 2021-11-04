package shell

import (
	"fmt"
	"testing"
)

type TestAlwaysFailCommand struct {
	Command
}

func (tafc TestAlwaysFailCommand) Run() error {
	return fmt.Errorf("always fail")
}

func Test_RunCommand(t *testing.T) {
	tafc := TestAlwaysFailCommand{}
	err := tafc.Run()
	if err == nil {
		t.Fatalf("expected to fail")
	}

	empty := NewCommand([]string{})
	err = empty.Run()
	if err == nil {
		t.Fatalf("empty command should error")
	}

	invalid := NewCommand([]string{"fdgdfgvcxvbxcvb", "cxzbxchthgf"})
	err = invalid.Run()
	if err == nil {
		t.Fatalf("invalid command should error")
	}

	invalid2 := NewCommand([]string{"fdgdfgvcxvbxcvb", "cxzbxchthgf"})
	_, err = invalid2.GetStdout()
	if err == nil {
		t.Fatalf("invalid command should error on getting output")
	}

	ls := NewCommand([]string{"ls"})
	out, err := ls.GetStdout()
	if err != nil {
		t.Log(err)
		t.Fatal("ls command should succeed")
	}

	if len(out) <= 0 {
		t.Fatalf("expected ls to return files list")
	}
}

func Test_GetShellBinaryDefaultsToBinSh(t *testing.T) {
	bin := GetShellBinary()
	if "/bin/sh" != bin {
		t.Fatalf("expected /bin/sh as default shell binary, got %s", bin)
	}
}
