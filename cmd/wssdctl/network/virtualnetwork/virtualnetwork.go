// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetwork

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetwork/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetwork",
		Aliases: []string{"vnet"},
		Short:   "network vnet resource",
		Long:    "network vnet resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
