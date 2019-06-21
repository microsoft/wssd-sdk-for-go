// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/create/virtualmachine"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "",
		Long:  "",
	}

	cmd.AddCommand(create.NewCommand())

	return cmd
}
