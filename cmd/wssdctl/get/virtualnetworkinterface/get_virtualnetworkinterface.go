// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetworkinterface

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetworkinterface",
		Aliases: []string{"vnic"},
		Short:   "Get a specific/all Virtual Network Interface(s)",
		Long:    "Get a specific/all Virtual Network Interface(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource")

	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")
	vnetInterfaceClient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	networkInterfaces, err := vnetInterfaceClient.Get(ctx, flags.Name)
	if err != nil {
		return err
	}
	// If a single vNET interface was requested
	if len(flags.Name) > 0 {
		if networkInterfaces == nil || len(*networkInterfaces) == 0 {
			return fmt.Errorf("Unable to find Virtual Network Interface [%s]", flags.Name)
		}

	} else {
		if networkInterfaces == nil || len(*networkInterfaces) == 0 {
			fmt.Println("No Virtual Network Interface Resources")
			// Not an error
			return nil
		}
	}

	virtualnetworkinterface.PrintList(networkInterfaces)

	return nil
}
