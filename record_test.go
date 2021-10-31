package main

import "testing"

func Test_RecordStore(t *testing.T) {
	store := newRecordStore("tmp/test.json")
	store.recordTime("test")

	tm := store.getTime("test")
	if tm <= 0 {
		t.Fatalf("recorded time should be larger than 0")
	}

	err := store.save()
	if err != nil {
		t.Log(err)
		t.Fatalf("saving should be a success")
	}

	store2 := newRecordStore("tmp/test.json")
	if store2.getTime("test") != store.getTime("test") {
		t.Fatal("expected stored time to match live time")
	}
}
