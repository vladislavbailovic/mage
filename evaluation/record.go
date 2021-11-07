package evaluation

import (
	"encoding/json"
	"io/ioutil"
	"mage/debug"
	"mage/epoch"
	"mage/typedefs"
	"os"
)

type recordStore struct {
	path    string
	records map[string]typedefs.Record
}

func NewRecordStore(path string) *recordStore {
	file, err := os.Open(path)
	if err != nil {
		return &recordStore{path, map[string]typedefs.Record{}}
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var records map[string]typedefs.Record
	json.Unmarshal(bytes, &records)

	store := recordStore{path, records}
	return &store
}

func (rs *recordStore) recordTime(entry string) {
	r := typedefs.Record{entry, epoch.Now()}
	rs.records[entry] = r
}

func (rs recordStore) getTime(entry string) typedefs.Epoch {
	r, ok := rs.records[entry]
	if !ok {
		return typedefs.Epoch(0)
	}
	return r.Timestamp
}

func (rs recordStore) save() error {
	debug.Records(rs.records)
	content, err := json.MarshalIndent(rs.records, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(rs.path, content, 0644)
	return err
}
