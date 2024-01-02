package utils

import (
	"encoding/json"
	"fmt"
)

func Print(obj any) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(data))
	return err
}
