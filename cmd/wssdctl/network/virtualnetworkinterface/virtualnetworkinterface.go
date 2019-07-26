// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetworkinterface

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network/virtualnetworkinterface/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualnetworkinterface",
		Aliases: []string{"vnic"},
		Short:   "network vnic resource",
		Long:    "network vnic resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
