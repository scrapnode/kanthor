package account

import (
	"encoding/json"
	"fmt"
	"os"
)

type output struct {
	stdout []string
	json   map[string]any
}

func (o *output) AddStdout(out string) {
	o.stdout = append(o.stdout, out)
}

func (o *output) RenderStdout() {
	for _, out := range o.stdout {
		fmt.Println(out)
	}
}

func (o *output) AddJson(name string, out any) {
	o.json[name] = out
}

func (o *output) RenderJson(dest string) error {
	bytes, err := json.MarshalIndent(o.json, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytes, os.ModePerm)
}

func (o *output) Render(dest string) error {
	if dest == "" {
		o.RenderStdout()
		return nil
	}

	return o.RenderJson(dest)
}
