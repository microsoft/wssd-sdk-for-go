// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package container

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/container/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/container/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/container/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/container/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage/container/update"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "container",
		Short: "container resource",
		Long:  "container resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
