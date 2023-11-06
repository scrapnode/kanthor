package account

import (
	"encoding/json"
	"fmt"
	"os"
)

type printing struct {
	stdout []string
	json   map[string]any
}

func (o *printing) AddStdout(out string) {
	o.stdout = append(o.stdout, out)
}

func (o *printing) RenderStdout() {
	for _, out := range o.stdout {
		fmt.Println(out)
	}
}

func (o *printing) AddJson(name string, out any) {
	o.json[name] = out
}

func (o *printing) RenderJson(dest string) error {
	bytes, err := json.MarshalIndent(o.json, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytes, os.ModePerm)
}

func (o *printing) Render(dest string) error {
	if dest == "" {
		o.RenderStdout()
		return nil
	}

	return o.RenderJson(dest)
}
