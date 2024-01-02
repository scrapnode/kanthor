package create

import (
	"context"
	"encoding/json"

	kanthor "github.com/scrapnode/kanthor/clients/sdk-go"
	"github.com/scrapnode/kanthor/cmd/client/utils"
	"github.com/spf13/cobra"
)

func Message() *cobra.Command {
	command := &cobra.Command{
		Use: "message",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := utils.ClientFromCmd(cmd)
			if err != nil {
				return err
			}
			req, err := parseMessageCreateReq(cmd)
			if err != nil {
				return err
			}
			res, err := client.Message.Create(context.Background(), req)
			if err != nil {
				return err
			}

			return utils.ResponseTo(cmd, res)
		},
	}

	command.PersistentFlags().StringP("app-id", "", "", "--app-id | the appplication that you want to make a message")
	command.PersistentFlags().StringP("msg-type", "", "", "--msg-type | the type of making message")
	command.PersistentFlags().StringP("msg-body", "", "", "--msg-body | the body of making message")

	return command
}

func parseMessageCreateReq(cmd *cobra.Command) (*kanthor.MessageCreateReq, error) {
	req := &kanthor.MessageCreateReq{}
	utils.RequestFrom(cmd, req)

	appId, err := cmd.Flags().GetString("app-id")
	if err != nil {
		return nil, err
	}
	if appId != "" {
		req.SetAppId(appId)
	}

	msgType, err := cmd.Flags().GetString("msg-type")
	if err != nil {
		return nil, err
	}
	if msgType != "" {
		req.SetType(msgType)
	}

	msgBody, err := cmd.Flags().GetString("msg-body")
	if err != nil {
		return nil, err
	}
	if msgBody != "" {
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(msgBody), &body); err != nil {
			return nil, err
		}
		req.SetBody(body)
	}

	return req, nil
}
