// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package identity

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/identity/create"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
	Group        string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "identity",
		Short: "identity resource",
		Long:  "identity resource",
	}

	cmd.AddCommand(create.NewCommand())

	return cmd
}
