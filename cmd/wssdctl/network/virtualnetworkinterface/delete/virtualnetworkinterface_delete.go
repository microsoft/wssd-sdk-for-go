// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Delete a specific Virtual Network Interface",
		Long:    "Delete a specific Virtual Network Interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network interface resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")
	vnetInterfaceClient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vnicName := flags.Name
	vnicId := ""
	if len(vnicName) == 0 {
		config := viper.GetString("config")
		vmconfig, err := virtualnetworkinterface.LoadConfig(config)
		if err != nil {
			return err
		}
		vnicName = *(vmconfig.Name)
		vnicId = *(vmconfig.ID)
	}

	err = vnetInterfaceClient.Delete(ctx, vnicName, vnicId)
	if err != nil {
		return err
	}

	return nil
}
