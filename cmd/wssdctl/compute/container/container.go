// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package container

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/list_ip_addresses"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/open_port"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/start"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/stop"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container/update"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "container",
		Short: "container compute resource",
		Long:  "container compute resource",
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
