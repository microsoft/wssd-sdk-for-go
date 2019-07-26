// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/loadbalancer/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "loadbalancer",
		Aliases: []string{"lb"},
		Short:   "network lb resource",
		Long:    "network lb resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
