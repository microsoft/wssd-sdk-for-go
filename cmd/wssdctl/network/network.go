// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package network

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "network",
		Short: "network resource",
		Long:  "network resource",
	}

	cmd.AddCommand(loadbalancer.NewCommand())
	cmd.AddCommand(virtualnetwork.NewCommand())
	cmd.AddCommand(virtualnetworkinterface.NewCommand())

	return cmd
}
