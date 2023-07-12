package cache

import "encoding/json"

func Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func Unmarshal[T any](data []byte) (T, error) {
	var dest T
	if err := json.Unmarshal(data, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}
