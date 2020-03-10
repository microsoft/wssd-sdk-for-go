// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package logging

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/logging/get"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "log",
		Short: "log operations",
		Long:  "log operations",
	}

	cmd.AddCommand(get.NewCommand())

	return cmd
}
