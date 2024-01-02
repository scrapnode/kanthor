package utils

import (
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
	credentials, err := cmd.Flags().GetString("credentials")
	if err != nil {
		return nil, err
	}
	if credentials == "" {
		if creds := os.Getenv("KANTHOR_CLIENT_CREDENTIALS"); creds != "" {
			credentials = creds
		}
	}

	opts := &kanthor.Options{Debug: debug, Host: host}
	return kanthor.NewWithOptions(credentials, opts)
}
