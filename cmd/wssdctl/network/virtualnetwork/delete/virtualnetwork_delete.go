// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

	server := viper.GetString("server")
	vnetClient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vnetName := flags.Name
	vnetId := ""
	if len(vnetName) == 0 {
		config := viper.GetString("config")
		vmconfig, err := virtualnetwork.LoadConfig(config)
		if err != nil {
			return err
		}
		vnetName = *(vmconfig.Name)
		vnetId = *(vmconfig.ID)
	}

	err = vnetClient.Delete(ctx, vnetName, vnetId)
	if err != nil {
		return err
	}

	return nil
}
