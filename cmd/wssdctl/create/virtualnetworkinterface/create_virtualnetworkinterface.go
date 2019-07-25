// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetworkinterface

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetworkinterface",
		Aliases: []string{"vnic"},
		Short:   "Create a Virtual Network interface",
		Long:    "Create a Virtual Network interface",
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
	vnicclient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server)
	if err != nil {
		return err
	}

	config := flags.FilePath
	vnicconfig, err := virtualnetworkinterface.LoadConfig(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = vnicclient.CreateOrUpdate(ctx, *(vnicconfig.Name), *(vnicconfig.ID), vnicconfig)
	if err != nil {
		return err
	}

	return nil
}
