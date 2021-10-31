package main

import (
	"os"
	"time"
	"io/ioutil"
	"encoding/json"
)

type epoch int64

type record struct {
	Name string
	Timestamp epoch
}

type recordStore struct {
	path string
	records map[string]record
}

func newRecordStore(path string) *recordStore {
	file, err := os.Open(path)
	if err != nil {
		return &recordStore{path, map[string]record{}}
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var records map[string]record
	json.Unmarshal(bytes, &records)

	store := recordStore{ path, records }
	return &store
}

func (rs *recordStore)recordTime(entry string) {
	r := record{entry, epoch(time.Now().Unix()) }
	rs.records[entry] = r
}

func (rs recordStore)getTime(entry string) epoch {
	r, ok := rs.records[entry]
	if !ok {
		return epoch(0)
	}
	return r.Timestamp
}

func (rs recordStore)save() error {
	content, err := json.MarshalIndent(rs.records, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(rs.path, content, 0644)
	return err
}
