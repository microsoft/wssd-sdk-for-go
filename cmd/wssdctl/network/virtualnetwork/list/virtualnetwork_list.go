// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"

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
		Use:   "list",
		Short: "Get all Virtual Network(s)",
		Long:  "Get all Virtual Network(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(flags *flags) error {

	server := viper.GetString("server")
	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networks, err := vnetclient.Get(ctx, flags.Name)
	if err != nil {
		return err
	}
	if networks == nil || len(*networks) == 0 {
		fmt.Println("No Virtual Network Resources")
		// Not an error
		return nil
	}

	virtualnetwork.PrintList(networks)

	return nil
}
