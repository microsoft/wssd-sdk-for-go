// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetwork",
		Aliases: []string{"vnet"},
		Short:   "Get a specific/all Virtual Network(s)",
		Long:    "Get a specific/all Virtual Network(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource")

	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")
	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	networks, err := vnetclient.Get(ctx, flags.Name)
	if err != nil {
		return err
	}
	// If a single vNET was requested
	if len(flags.Name) > 0 {
		if networks == nil || len(*networks) == 0 {
			return fmt.Errorf("Unable to find Virtual Network [%s]", flags.Name)
		}

	} else {
		if networks == nil || len(*networks) == 0 {
			fmt.Println("No Virtual Network Resources")
			// Not an error
			return nil
		}
	}

	virtualnetwork.PrintList(networks)

	return nil
}
