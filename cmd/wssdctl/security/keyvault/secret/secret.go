// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package secret

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/keyvault/secret/set"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/keyvault/secret/show"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/keyvault/secret/download"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "secret",
		Short: "secret resource",
		Long:  "secret resource",
	}

	cmd.AddCommand(set.NewCommand())
	cmd.AddCommand(show.NewCommand())
	cmd.AddCommand(download.NewCommand())
	
	return cmd
}
