// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create a vhd",
		Long:  "Create a vhd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")
	container := viper.GetString("container")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vhdClient, err := virtualharddisk.NewVirtualHardDiskClient(server, authorizer)
	if err != nil {
		return err
	}

	vhdConfig := storage.VirtualHardDisk{}
	err = config.LoadYAMLFile(flags.FilePath, &vhdConfig)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if vhdConfig.Name == nil {
		return errors.Wrapf(errors.InvalidInput, "Missing Name")
	}

	_, err = vhdClient.CreateOrUpdate(ctx, container, *(vhdConfig.Name), &vhdConfig)
	if err != nil {
		return err
	}

	return nil
}
