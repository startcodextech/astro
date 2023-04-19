package utils

import (
	"github.com/goccy/go-json"
	"os"
	"path/filepath"
)

func JsonExist(filename string, isArray bool) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {

		err := os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return err
		}

		var data []byte
		if isArray {
			emptyArray := make([]interface{}, 0)
			data, err = json.Marshal(emptyArray)
			if err != nil {
				return err
			}
		} else {
			emptyObject := make(map[string]interface{})
			data, err = json.Marshal(emptyObject)
			if err != nil {
				return err
			}
		}

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func LoadJson(path string, out interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(out)
	if err != nil {
		return err
	}

	return nil
}
