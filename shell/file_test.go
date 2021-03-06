package shell

import "testing"

func Test_FileExists(t *testing.T) {
	if FileExists("whatever, this does not exist") {
		t.Fatalf("non-existent file should not exist")
	}
	if !FileExists("../fixtures/simple.mg") {
		t.Fatalf("valid file should exist")
	}
}

func Test_GetFileMtime(t *testing.T) {
	mtime := GetFileMtime("whatever no such file")
	if 0 != mtime {
		t.Fatalf("mtime of invalid file should be 0, got %d", mtime)
	}

	mtime = GetFileMtime("../fixtures/simple.mg")
	if mtime <= 0 {
		t.Fatalf("mtime of valid file should NOT be 0")
	}
}

func Test_LoadFile(t *testing.T) {
	_, err := LoadFile("no such file")
	if err == nil {
		t.Fatalf("should error out when loading invalid file")
	}

	stuff, _ := LoadFile("../fixtures/simple.mg")
	if len(stuff) <= 0 {
		t.Fatalf("should have loaded valid file")
	}
}
