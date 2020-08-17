package Core

import "encoding/json"

func remapJson(source interface{}, destination interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, destination)
}
