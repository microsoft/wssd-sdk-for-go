// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package get

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/all"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/virtualmachine"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "get",
		Short: "Get a resource",
		Long:  "Get a resource",
	}
	cmd.AddCommand(all.NewCommand())
	cmd.AddCommand(virtualmachine.NewCommand())

	return cmd
}
