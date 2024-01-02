package utils

import (
	"encoding/json"
	"os"

	kanthor "github.com/scrapnode/kanthor/clients/sdk-go"
	"github.com/spf13/cobra"
)

func ClientFromCmd(cmd *cobra.Command) (*kanthor.Kanthor, error) {
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil, err
	}
	host, err := cmd.Flags().GetString("host")
	if err != nil {
		return nil, err
	}
	if host == "" {
		if h := os.Getenv("KANTHOR_SERVER_HOST"); h != "" {
			host = h
		}
	}

	credentials, err := cmd.Flags().GetString("credentials")
	if err != nil {
		return nil, err
	}
	if credentials == "" {
		if creds := os.Getenv("KANTHOR_SERVER_CREDENTIALS"); creds != "" {
			credentials = creds
		}
	}

	opts := &kanthor.Options{Debug: debug, Host: host}
	return kanthor.NewWithOptions(credentials, opts)
}

func RequestFrom(cmd *cobra.Command, req any) error {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}
	if input == "" {
		return nil
	}

	data, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, req)
}

func ResponseTo(cmd *cobra.Command, res any) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if output == "" {
		return Print(res)
	}

	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(output, data, os.ModePerm)
}
