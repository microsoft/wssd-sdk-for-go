// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package compute

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/container"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachinescaleset"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "compute",
		Short: "compute resource",
		Long:  "compute resource",
	}

	cmd.AddCommand(container.NewCommand())
	cmd.AddCommand(virtualmachine.NewCommand())
	cmd.AddCommand(virtualmachinescaleset.NewCommand())

	return cmd
}
