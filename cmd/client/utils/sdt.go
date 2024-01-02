package utils

import (
	"encoding/json"
	"fmt"
)

func Print(obj any) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
