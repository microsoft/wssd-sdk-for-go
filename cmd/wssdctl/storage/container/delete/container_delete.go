// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"
	"time"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/storage/container"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "delete a specific container",
		Long:  "delete a specific container ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the container")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	containerClient, err := container.NewContainerClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = containerClient.Delete(ctx, group, flags.Name)
	if err != nil {
		return err
	}

	return nil
}
