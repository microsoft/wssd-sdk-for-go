// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"

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
		Args:  cobra.NoArgs,
		Use:   "list",
		Short: "Get a all Virtual Network Interface(s)",
		Long:  "Get a all Virtual Network Interface(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

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

	networkInterfaces, err := vnetInterfaceClient.Get(ctx, "")
	if err != nil {
		return err
	}

	if networkInterfaces == nil || len(*networkInterfaces) == 0 {
		fmt.Println("No Virtual Network Interface Resources")
		// Not an error
		return nil
	}

	virtualnetworkinterface.PrintList(networkInterfaces)

	return nil
}
