// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/delete/virtualmachine"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "delete a Virtual Machine resource",
		Long:  "delete a Virtual Machine resource",
	}
	cmd.AddCommand(virtualmachine.NewCommand())

	return cmd
}
