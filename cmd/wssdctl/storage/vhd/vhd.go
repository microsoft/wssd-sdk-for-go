// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package vhd

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/vhd/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/vhd/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/vhd/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/vhd/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/vhd/update"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "vhd",
		Short: "vhd resource",
		Long:  "vhd resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
