// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachinescaleset

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/list_ip_addresses"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/list_virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/start"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/stop"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachinescaleset",
		Aliases: []string{"vmss"},
		Short:   "virtualmachinescaleset compute resource",
		Long:    "virtualmachinescaleset compute resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(list_ip_addresses.NewCommand())
	cmd.AddCommand(list_virtualmachine.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(start.NewCommand())
	cmd.AddCommand(stop.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
