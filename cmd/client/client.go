package client

import (
	"github.com/scrapnode/kanthor/cmd/client/create"
	"github.com/spf13/cobra"
)

var example = `
kanthor client --host=localhost:8180 --debug create message
`

func New() *cobra.Command {
	command := &cobra.Command{
		Use:     "client",
		Example: example,
	}

	command.AddCommand(create.New())

	command.PersistentFlags().BoolP("debug", "", false, "--debug | show debug information when client perform an action")
	command.PersistentFlags().StringP("host", "", "", "--host | override destination sever host")
	command.PersistentFlags().StringP("credentials", "", "", "--credentials | the basic authentication in the form of <USER>:<PASS>")

	command.PersistentFlags().StringP("input", "i", "", "--input | the file that contains request payload of modification APIs (POST/PUT/PATCH)")
	command.PersistentFlags().StringP("output", "o", "", "--output | the output file that will be writen with API response")

	return command
}
