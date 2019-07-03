// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package resourcegroup

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "resourcegroup",
		Aliases: []string{"server", "machine"},
		Short:   "cloud vnet resource",
		Long:    "cloud vnet resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
