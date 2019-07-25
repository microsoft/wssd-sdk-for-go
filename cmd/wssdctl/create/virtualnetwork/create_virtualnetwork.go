// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetwork

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetwork",
		Aliases: []string{"vnet"},
		Short:   "Create a Virtual Network",
		Long:    "Create a Virtual Network",
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
	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return err
	}

	config := flags.FilePath
	vnetconfig, err := virtualnetwork.LoadConfig(config)
	if err != nil {
		return err
	}

	// Wait up to one minute for network creation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = vnetclient.CreateOrUpdate(ctx, *(vnetconfig.Name), *(vnetconfig.ID), vnetconfig)
	if err != nil {
		return err
	}

	return nil
}
