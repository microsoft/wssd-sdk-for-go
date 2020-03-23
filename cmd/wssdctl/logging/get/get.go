// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package get

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/admin/logging"
)

type flags struct {
	Outname string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "get",
		Short: "Get full log file from agent",
		Long:  "Get full log file from agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Outname, "outname", "agent.log", "name to output log to")
	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	logClient, err := logging.NewLoggingClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	err = logClient.GetLogFile(ctx, flags.Outname)

	return err
}
