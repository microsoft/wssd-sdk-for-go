// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachine

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/list_ip_addresses"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/open_port"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/start"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/stop"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachine",
		Aliases: []string{"vm"},
		Short:   "virtualmachine compute resource",
		Long:    "virtualmachine compute resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(list_ip_addresses.NewCommand())
	cmd.AddCommand(open_port.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(start.NewCommand())
	cmd.AddCommand(stop.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
