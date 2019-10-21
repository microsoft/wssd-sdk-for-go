// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "Delete a specific/all Virtual Network(s)",
		Long:  "Delete a specific/all Virtual Network(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	vnetClient, err := virtualnetwork.NewVirtualNetworkClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vnetName := flags.Name
	if len(vnetName) == 0 {
		vnetconfig := network.VirtualNetwork{}
		configFile := viper.GetString("config")
		err = config.LoadYAMLFile(configFile, &vnetconfig)
		if err != nil {
			return err
		}
		vnetName = *(vnetconfig.Name)
	}

	err = vnetClient.Delete(ctx, group, vnetName)
	if err != nil {
		return err
	}

	return nil
}
