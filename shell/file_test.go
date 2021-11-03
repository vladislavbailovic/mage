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
