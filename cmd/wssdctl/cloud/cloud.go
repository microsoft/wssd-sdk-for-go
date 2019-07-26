// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package cloud

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/node"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/resourcegroup"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/cloud/update"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cloud",
		Short: "cloud resource",
		Long:  "cloud resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(node.NewCommand())
	cmd.AddCommand(resourcegroup.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(update.NewCommand())

	return cmd
}
