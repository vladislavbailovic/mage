package ruleset

import (
	"testing"
)

func Test_RecordStore(t *testing.T) {
	store := NewRecordStore("../tmp/test.json")
	store.RecordTime("test")

	tm := store.GetTime("test")
	if tm <= 0 {
		t.Fatalf("recorded time should be larger than 0")
	}

	err := store.Save()
	if err != nil {
		t.Log(err)
		t.Fatalf("saving should be a success")
	}

	store2 := NewRecordStore("../tmp/test.json")
	if store2.GetTime("test") != store.GetTime("test") {
		t.Fatal("expected stored time to match live time")
	}
}
