package encodingutil

import (
	"encoding/json"
	"io/ioutil"
)

func MarshalJSONToFile(fn string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fn, data, 0666)
}

func UnmarshalJSONFromFile(fn string, v interface{}) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
