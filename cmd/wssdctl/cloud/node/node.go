// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package node

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "node",
		Aliases: []string{"server", "machine"},
		Short:   "cloud node resource",
		Long:    "cloud node resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
